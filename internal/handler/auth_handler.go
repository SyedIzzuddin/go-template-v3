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