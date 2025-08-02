package middleware

import (
	"go-template/internal/logger"
	"go-template/pkg/jwt"
	"go-template/pkg/response"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(jwtManager *jwt.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header", zap.String("request_id", requestID))
				return response.Unauthorized(c, "Authorization header required")
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Warn("Invalid Authorization header format", zap.String("request_id", requestID))
				return response.Unauthorized(c, "Authorization header must start with 'Bearer '")
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				logger.Warn("Empty token in Authorization header", zap.String("request_id", requestID))
				return response.Unauthorized(c, "Token cannot be empty")
			}

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(token)
			if err != nil {
				logger.Warn("Invalid or expired token", 
					zap.Error(err), 
					zap.String("request_id", requestID))
				
				switch err {
				case jwt.ErrExpiredToken:
					return response.Unauthorized(c, "Token has expired")
				case jwt.ErrInvalidToken, jwt.ErrInvalidClaims:
					return response.Unauthorized(c, "Invalid token")
				default:
					return response.Unauthorized(c, "Token validation failed")
				}
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			logger.Debug("Authentication successful", 
				zap.String("request_id", requestID),
				zap.Int("user_id", claims.UserID),
				zap.String("user_email", claims.Email))

			return next(c)
		}
	}
}

// OptionalAuthMiddleware creates optional JWT authentication middleware
// Sets user info in context if valid token is provided, but doesn't require it
func OptionalAuthMiddleware(jwtManager *jwt.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without authentication
				return next(c)
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// Invalid format, continue without authentication
				return next(c)
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				// Empty token, continue without authentication
				return next(c)
			}

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(token)
			if err != nil {
				logger.Debug("Optional auth failed", 
					zap.Error(err), 
					zap.String("request_id", requestID))
				// Invalid token, continue without authentication
				return next(c)
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			logger.Debug("Optional authentication successful", 
				zap.String("request_id", requestID),
				zap.Int("user_id", claims.UserID),
				zap.String("user_email", claims.Email))

			return next(c)
		}
	}
}