package service

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/config"
	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/repository"
	"go-template/pkg/email"
	"go-template/pkg/jwt"
	"go-template/pkg/tokens"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials         = errors.New("invalid email or password")
	ErrUserAlreadyExists          = errors.New("user with this email already exists")
	ErrInvalidVerificationToken   = errors.New("invalid or expired verification token")
	ErrEmailAlreadyVerified       = errors.New("email is already verified")
	ErrInvalidPasswordResetToken  = errors.New("invalid or expired password reset token")
	ErrUserNotFound               = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	GetUserProfile(ctx context.Context, userID int) (*dto.UserProfileResponse, error)
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (*dto.EmailVerificationResponse, error)
	ResendVerificationEmail(ctx context.Context, req dto.ResendVerificationRequest) (*dto.EmailVerificationResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.PasswordResetResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.PasswordResetResponse, error)
}

type authService struct {
	userRepo     repository.UserRepository
	jwtManager   *jwt.JWTManager
	emailService email.Service
	config       *config.Config
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager, emailService email.Service, config *config.Config) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		emailService: emailService,
		config:       config,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	logger.Info("User registration attempt", zap.String("email", req.Email))

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		logger.Warn("Registration attempt with existing email", zap.String("email", req.Email))
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("failed to process password")
	}

	// Generate email verification token
	verificationToken, err := tokens.GenerateVerificationToken()
	if err != nil {
		logger.Error("Failed to generate verification token", zap.Error(err))
		return nil, errors.New("failed to generate verification token")
	}

	// Set token expiration (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create user with hashed password and verification token
	user, err := s.userRepo.CreateWithPasswordAndRole(ctx, req.Name, req.Email, hashedPassword, "user", verificationToken, &expiresAt)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Email, user.Name, verificationToken); err != nil {
		logger.Error("Failed to send verification email", zap.Error(err))
		// Don't fail registration if email sending fails - just log it
		logger.Warn("User registered but verification email not sent", zap.Int("user_id", user.ID))
	} else {
		logger.Info("Verification email sent", zap.Int("user_id", user.ID))
	}

	// Generate token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, errors.New("failed to generate authentication tokens")
	}

	logger.Info("User registered successfully", zap.Int("user_id", user.ID))

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Role:          user.Role,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessExpiresIn),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	logger.Info("User login attempt", zap.String("email", req.Email))

	// Get user by email
	user, err := s.userRepo.GetByEmailWithPassword(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Login attempt with non-existent email", zap.String("email", req.Email))
			return nil, ErrInvalidCredentials
		}
		logger.Error("Failed to get user", zap.Error(err))
		return nil, errors.New("authentication failed")
	}

	// Verify password
	if err := s.verifyPassword(req.Password, user.PasswordHash); err != nil {
		logger.Warn("Login attempt with invalid password", zap.String("email", req.Email))
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, errors.New("failed to generate authentication tokens")
	}

	logger.Info("User logged in successfully", zap.Int("user_id", user.ID))

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Role:          user.Role,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.AccessExpiresIn),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	logger.Info("Token refresh attempt")

	// Generate new access token using refresh token
	newAccessToken, err := s.jwtManager.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		logger.Warn("Failed to refresh token", zap.Error(err))
		return nil, errors.New("invalid refresh token")
	}

	logger.Info("Token refreshed successfully")

	return &dto.TokenResponse{
		AccessToken: newAccessToken,
		ExpiresAt:   time.Now().Add(s.config.JWT.AccessExpiresIn),
	}, nil
}

func (s *authService) GetUserProfile(ctx context.Context, userID int) (*dto.UserProfileResponse, error) {
	logger.Debug("Getting user profile", zap.Int("user_id", userID))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Profile request for non-existent user", zap.Int("user_id", userID))
			return nil, errors.New("user not found")
		}
		logger.Error("Failed to get user profile", zap.Error(err))
		return nil, errors.New("failed to get user profile")
	}

	return &dto.UserProfileResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

// hashPassword hashes a plain text password using bcrypt
func (s *authService) hashPassword(password string) (string, error) {
	// Use bcrypt with cost 12 for strong security
	// Cost 12 is recommended by OWASP for 2025
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword compares a plain text password with a hashed password
func (s *authService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *authService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (*dto.EmailVerificationResponse, error) {
	logger.Info("Email verification attempt", zap.String("token", req.Token[:8]+"..."))

	// Get user by verification token
	user, err := s.userRepo.GetByVerificationToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Invalid verification token used", zap.String("token", req.Token[:8]+"..."))
			return &dto.EmailVerificationResponse{
				Message: "Invalid or expired verification token",
				Success: false,
			}, ErrInvalidVerificationToken
		}
		logger.Error("Failed to get user by verification token", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Verification failed",
			Success: false,
		}, errors.New("verification failed")
	}

	// Check if email is already verified
	if user.EmailVerified {
		logger.Info("Attempt to verify already verified email", zap.Int("user_id", user.ID))
		return &dto.EmailVerificationResponse{
			Message: "Email is already verified",
			Success: false,
		}, ErrEmailAlreadyVerified
	}

	// Check if token has expired
	if user.EmailVerificationExpiresAt != nil && time.Now().After(*user.EmailVerificationExpiresAt) {
		logger.Warn("Expired verification token used", zap.Int("user_id", user.ID))
		return &dto.EmailVerificationResponse{
			Message: "Verification token has expired",
			Success: false,
		}, ErrInvalidVerificationToken
	}

	// Verify the email
	if err := s.userRepo.VerifyEmail(ctx, req.Token); err != nil {
		logger.Error("Failed to verify email", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Verification failed",
			Success: false,
		}, errors.New("verification failed")
	}

	logger.Info("Email verified successfully", zap.Int("user_id", user.ID))

	return &dto.EmailVerificationResponse{
		Message: "Email verified successfully",
		Success: true,
	}, nil
}

func (s *authService) ResendVerificationEmail(ctx context.Context, req dto.ResendVerificationRequest) (*dto.EmailVerificationResponse, error) {
	logger.Info("Resend verification email attempt", zap.String("email", req.Email))

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Resend verification attempt for non-existent email", zap.String("email", req.Email))
			return &dto.EmailVerificationResponse{
				Message: "User not found",
				Success: false,
			}, errors.New("user not found")
		}
		logger.Error("Failed to get user by email", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Failed to resend verification email",
			Success: false,
		}, errors.New("failed to resend verification email")
	}

	// Check if email is already verified
	if user.EmailVerified {
		logger.Info("Resend verification attempt for already verified email", zap.Int("user_id", user.ID))
		return &dto.EmailVerificationResponse{
			Message: "Email is already verified",
			Success: false,
		}, ErrEmailAlreadyVerified
	}

	// Generate new verification token
	verificationToken, err := tokens.GenerateVerificationToken()
	if err != nil {
		logger.Error("Failed to generate verification token", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Failed to generate verification token",
			Success: false,
		}, errors.New("failed to generate verification token")
	}

	// Set new token expiration (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Update user with new verification token
	if err := s.userRepo.UpdateVerificationToken(ctx, user.ID, verificationToken, &expiresAt); err != nil {
		logger.Error("Failed to update verification token", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Failed to update verification token",
			Success: false,
		}, errors.New("failed to update verification token")
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Email, user.Name, verificationToken); err != nil {
		logger.Error("Failed to send verification email", zap.Error(err))
		return &dto.EmailVerificationResponse{
			Message: "Failed to send verification email",
			Success: false,
		}, errors.New("failed to send verification email")
	}

	logger.Info("Verification email resent successfully", zap.Int("user_id", user.ID))

	return &dto.EmailVerificationResponse{
		Message: "Verification email sent successfully",
		Success: true,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (*dto.PasswordResetResponse, error) {
	logger.Info("Forgot password attempt", zap.String("email", req.Email))

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Forgot password attempt for non-existent email", zap.String("email", req.Email))
			// For security, don't reveal that the email doesn't exist
			return &dto.PasswordResetResponse{
				Message: "If your email is registered, you will receive a password reset link shortly",
				Success: true,
			}, nil
		}
		logger.Error("Failed to get user by email", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to process password reset request",
			Success: false,
		}, errors.New("failed to process password reset request")
	}

	// Generate password reset token
	resetToken, err := tokens.GeneratePasswordResetToken()
	if err != nil {
		logger.Error("Failed to generate password reset token", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to generate password reset token",
			Success: false,
		}, errors.New("failed to generate password reset token")
	}

	// Set token expiration (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Update user with password reset token
	if err := s.userRepo.UpdatePasswordResetToken(ctx, user.ID, resetToken, &expiresAt); err != nil {
		logger.Error("Failed to update password reset token", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to update password reset token",
			Success: false,
		}, errors.New("failed to update password reset token")
	}

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.Name, resetToken); err != nil {
		logger.Error("Failed to send password reset email", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to send password reset email",
			Success: false,
		}, errors.New("failed to send password reset email")
	}

	logger.Info("Password reset email sent successfully", zap.Int("user_id", user.ID))

	return &dto.PasswordResetResponse{
		Message: "If your email is registered, you will receive a password reset link shortly",
		Success: true,
	}, nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (*dto.PasswordResetResponse, error) {
	logger.Info("Password reset attempt", zap.String("token", req.Token[:8]+"..."))

	// Get user by password reset token
	user, err := s.userRepo.GetByPasswordResetToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Invalid password reset token used", zap.String("token", req.Token[:8]+"..."))
			return &dto.PasswordResetResponse{
				Message: "Invalid or expired password reset token",
				Success: false,
			}, ErrInvalidPasswordResetToken
		}
		logger.Error("Failed to get user by password reset token", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Password reset failed",
			Success: false,
		}, errors.New("password reset failed")
	}

	// Check if token has expired
	if user.PasswordResetExpiresAt != nil && time.Now().After(*user.PasswordResetExpiresAt) {
		logger.Warn("Expired password reset token used", zap.Int("user_id", user.ID))
		return &dto.PasswordResetResponse{
			Message: "Password reset token has expired",
			Success: false,
		}, ErrInvalidPasswordResetToken
	}

	// Hash the new password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash new password", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to process new password",
			Success: false,
		}, errors.New("failed to process new password")
	}

	// Reset the password
	if err := s.userRepo.ResetPassword(ctx, req.Token, hashedPassword); err != nil {
		logger.Error("Failed to reset password", zap.Error(err))
		return &dto.PasswordResetResponse{
			Message: "Failed to reset password",
			Success: false,
		}, errors.New("failed to reset password")
	}

	logger.Info("Password reset successfully", zap.Int("user_id", user.ID))

	return &dto.PasswordResetResponse{
		Message: "Password reset successfully",
		Success: true,
	}, nil
}