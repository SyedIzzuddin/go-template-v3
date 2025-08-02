package middleware

import (
	"context"
	"database/sql"
	"go-template/internal/logger"
	"go-template/internal/repository"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// EmailVerificationMiddleware adds a warning header for unverified users
// This implements the "soft block" approach where unverified users can still access the system
// but get notified about their unverified status
func EmailVerificationMiddleware(userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Check if user is authenticated (set by AuthMiddleware)
			userID, ok := c.Get("user_id").(int)
			if !ok {
				// No authenticated user, continue without verification check
				return next(c)
			}

			// Get user details to check email verification status
			user, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				if err != sql.ErrNoRows {
					logger.Warn("Failed to get user for email verification check", 
						zap.Error(err), 
						zap.String("request_id", requestID),
						zap.Int("user_id", userID))
				}
				// Continue without verification check if user lookup fails
				return next(c)
			}

			// Set email verification status in context for handlers to use
			c.Set("email_verified", user.EmailVerified)

			// Add warning header if email is not verified
			if !user.EmailVerified {
				c.Response().Header().Set("X-Email-Verification-Status", "unverified")
				c.Response().Header().Set("X-Email-Verification-Warning", "Please verify your email address to ensure full account security")
				
				logger.Debug("Unverified user accessing protected resource", 
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.String("user_email", user.Email))
			}

			return next(c)
		}
	}
}