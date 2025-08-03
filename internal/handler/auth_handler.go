package handler

import (
	"go-template/internal/dto"
	"go-template/internal/logger"
	"go-template/internal/service"
	"go-template/pkg/response"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validator
}

func NewAuthHandler(authService service.AuthService, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password. Default role is 'user'.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} response.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure 400 {object} response.Response{error=[]dto.ValidationError} "Validation error"
// @Failure 409 {object} response.Response "User already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("User registration request started", zap.String("request_id", requestID))

	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind registration request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Registration validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	authResponse, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err), zap.String("request_id", requestID))
		if err == service.ErrUserAlreadyExists {
			return response.Conflict(c, "User with this email already exists", nil)
		}
		return response.BadRequest(c, "Registration failed", err.Error())
	}

	logger.Info("User registration completed successfully", 
		zap.String("request_id", requestID), 
		zap.Int("user_id", authResponse.User.ID))
	return response.Created(c, "User registered successfully", authResponse)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} response.Response{data=dto.AuthResponse} "Login successful"
// @Failure 400 {object} response.Response{error=[]dto.ValidationError} "Validation error"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("User login request started", zap.String("request_id", requestID))

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind login request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Login validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	authResponse, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to login user", zap.Error(err), zap.String("request_id", requestID))
		if err == service.ErrInvalidCredentials {
			return response.Unauthorized(c, "Invalid email or password")
		}
		return response.InternalServerError(c, "Login failed", err.Error())
	}

	logger.Info("User login completed successfully", 
		zap.String("request_id", requestID), 
		zap.Int("user_id", authResponse.User.ID))
	return response.Success(c, "Login successful", authResponse)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("Token refresh request started", zap.String("request_id", requestID))

	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind refresh token request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Refresh token validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	tokenResponse, err := h.authService.RefreshToken(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to refresh token", zap.Error(err), zap.String("request_id", requestID))
		return response.Unauthorized(c, "Invalid refresh token")
	}

	logger.Info("Token refresh completed successfully", zap.String("request_id", requestID))
	return response.Success(c, "Token refreshed successfully", tokenResponse)
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserProfileResponse} "Profile retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	userID := c.Get("user_id").(int) // Set by auth middleware
	logger.Info("Get user profile request started", 
		zap.String("request_id", requestID), 
		zap.Int("user_id", userID))

	profile, err := h.authService.GetUserProfile(c.Request().Context(), userID)
	if err != nil {
		logger.Error("Failed to get user profile", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to get profile", err.Error())
	}

	logger.Info("Get user profile completed successfully", zap.String("request_id", requestID))
	return response.Success(c, "Profile retrieved successfully", profile)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("Email verification request started", zap.String("request_id", requestID))

	var req dto.VerifyEmailRequest
	
	// Try to get token from query parameter first (for email links)
	token := c.QueryParam("token")
	if token != "" {
		req.Token = token
	} else {
		// Otherwise, try to bind from request body (for API calls)
		if err := c.Bind(&req); err != nil {
			logger.Error("Failed to bind email verification request", zap.Error(err), zap.String("request_id", requestID))
			return response.BadRequest(c, "Invalid request body", err.Error())
		}
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Email verification validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	verifyResponse, err := h.authService.VerifyEmail(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to verify email", zap.Error(err), zap.String("request_id", requestID))
		if err == service.ErrInvalidVerificationToken {
			return response.BadRequest(c, verifyResponse.Message, nil)
		}
		if err == service.ErrEmailAlreadyVerified {
			return response.BadRequest(c, verifyResponse.Message, nil)
		}
		return response.InternalServerError(c, "Email verification failed", err.Error())
	}

	logger.Info("Email verification completed successfully", zap.String("request_id", requestID))
	return response.Success(c, verifyResponse.Message, verifyResponse)
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("Resend verification email request started", zap.String("request_id", requestID))

	var req dto.ResendVerificationRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind resend verification request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Resend verification validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	resendResponse, err := h.authService.ResendVerificationEmail(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to resend verification email", zap.Error(err), zap.String("request_id", requestID))
		if err == service.ErrEmailAlreadyVerified {
			return response.BadRequest(c, resendResponse.Message, nil)
		}
		return response.InternalServerError(c, "Failed to resend verification email", err.Error())
	}

	logger.Info("Resend verification email completed successfully", zap.String("request_id", requestID))
	return response.Success(c, resendResponse.Message, resendResponse)
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("Forgot password request started", zap.String("request_id", requestID))

	var req dto.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind forgot password request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Forgot password validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	resetResponse, err := h.authService.ForgotPassword(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to process forgot password", zap.Error(err), zap.String("request_id", requestID))
		return response.InternalServerError(c, "Failed to process password reset request", err.Error())
	}

	logger.Info("Forgot password completed successfully", zap.String("request_id", requestID))
	return response.Success(c, resetResponse.Message, resetResponse)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger.Info("Reset password request started", zap.String("request_id", requestID), zap.String("method", c.Request().Method))

	// Handle GET requests (from email links) - validate token and return instructions
	if c.Request().Method == "GET" {
		token := c.QueryParam("token")
		if token == "" {
			logger.Warn("GET request to reset password without token", zap.String("request_id", requestID))
			return response.BadRequest(c, "Password reset token is required", nil)
		}

		// Validate the token format
		if len(token) != 64 { // hex encoding of 32 bytes = 64 characters
			logger.Warn("Invalid token format in GET request", zap.String("request_id", requestID))
			return response.BadRequest(c, "Invalid password reset token format", nil)
		}

		logger.Info("Valid password reset token accessed via GET", zap.String("request_id", requestID))
		
		// Return instructions for the user
		return response.Success(c, "Password reset token is valid. Please use POST method with your new password to complete the reset.", map[string]interface{}{
			"token": token,
			"instructions": "Send a POST request to this same endpoint with 'token' and 'password' in the request body",
			"example": map[string]interface{}{
				"method": "POST",
				"body": map[string]interface{}{
					"token":    token,
					"password": "YourNewPassword123!",
				},
			},
		})
	}

	// Handle POST requests (actual password reset)
	var req dto.ResetPasswordRequest
	
	// Try to get token from query parameter first (for convenience)
	token := c.QueryParam("token")
	if token != "" {
		req.Token = token
	}

	// Bind from request body (will override token if provided in body)
	if err := c.Bind(&req); err != nil {
		logger.Error("Failed to bind reset password request", zap.Error(err), zap.String("request_id", requestID))
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// If token was from query param and not in body, use the query param
	if token != "" && req.Token == "" {
		req.Token = token
	}

	// Validate request
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		logger.Warn("Reset password validation failed", zap.Any("errors", validationErrors), zap.String("request_id", requestID))
		return response.ValidationError(c, "Validation failed", validationErrors)
	}

	resetResponse, err := h.authService.ResetPassword(c.Request().Context(), req)
	if err != nil {
		logger.Error("Failed to reset password", zap.Error(err), zap.String("request_id", requestID))
		if err == service.ErrInvalidPasswordResetToken {
			return response.BadRequest(c, resetResponse.Message, nil)
		}
		return response.InternalServerError(c, "Failed to reset password", err.Error())
	}

	logger.Info("Reset password completed successfully", zap.String("request_id", requestID))
	return response.Success(c, resetResponse.Message, resetResponse)
}