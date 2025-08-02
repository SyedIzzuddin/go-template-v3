package service

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/config"
	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/repository"
	"go-template/pkg/jwt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	GetUserProfile(ctx context.Context, userID int) (*dto.UserProfileResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
	config     *config.Config
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager, config *config.Config) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		config:     config,
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

	// Create user with hashed password
	user, err := s.userRepo.CreateWithPassword(ctx, req.Name, req.Email, hashedPassword)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
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
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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