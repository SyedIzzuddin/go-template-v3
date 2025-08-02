package middleware

import (
	"context"
	"database/sql"
	"go-template/internal/logger"
	"go-template/internal/repository"
	"go-template/pkg/response"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RoleMiddleware creates middleware that checks if user has required role
func RoleMiddleware(userRepo repository.UserRepository, requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get user ID from JWT token (set by AuthMiddleware)
			userID, ok := c.Get("user_id").(int)
			if !ok {
				logger.Warn("RBAC check failed: user not authenticated", 
					zap.String("request_id", requestID),
					zap.String("required_role", requiredRole))
				return response.Unauthorized(c, "User not authenticated")
			}

			// Get user from database to check role
			user, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				if err == sql.ErrNoRows {
					logger.Warn("RBAC check failed: user not found", 
						zap.String("request_id", requestID),
						zap.Int("user_id", userID),
						zap.String("required_role", requiredRole))
					return response.Unauthorized(c, "User not found")
				}
				
				logger.Error("RBAC check failed: database error", 
					zap.Error(err),
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.String("required_role", requiredRole))
				return response.InternalServerError(c, "Internal server error", nil)
			}

			// Check if user has required role
			if user.Role != requiredRole {
				logger.Warn("RBAC check failed: insufficient permissions", 
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.String("user_role", user.Role),
					zap.String("required_role", requiredRole),
					zap.String("user_email", user.Email))
				return response.Forbidden(c, "Insufficient permissions")
			}

			// Set user role in context for handlers to use
			c.Set("user_role", user.Role)

			logger.Debug("RBAC check passed", 
				zap.String("request_id", requestID),
				zap.Int("user_id", userID),
				zap.String("user_role", user.Role),
				zap.String("required_role", requiredRole))

			return next(c)
		}
	}
}

// MultiRoleMiddleware creates middleware that checks if user has any of the allowed roles
func MultiRoleMiddleware(userRepo repository.UserRepository, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get user ID from JWT token (set by AuthMiddleware)
			userID, ok := c.Get("user_id").(int)
			if !ok {
				logger.Warn("Multi-role RBAC check failed: user not authenticated", 
					zap.String("request_id", requestID),
					zap.Strings("allowed_roles", allowedRoles))
				return response.Unauthorized(c, "User not authenticated")
			}

			// Get user from database to check role
			user, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				if err == sql.ErrNoRows {
					logger.Warn("Multi-role RBAC check failed: user not found", 
						zap.String("request_id", requestID),
						zap.Int("user_id", userID),
						zap.Strings("allowed_roles", allowedRoles))
					return response.Unauthorized(c, "User not found")
				}
				
				logger.Error("Multi-role RBAC check failed: database error", 
					zap.Error(err),
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.Strings("allowed_roles", allowedRoles))
				return response.InternalServerError(c, "Internal server error", nil)
			}

			// Check if user has any of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				logger.Warn("Multi-role RBAC check failed: insufficient permissions", 
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.String("user_role", user.Role),
					zap.Strings("allowed_roles", allowedRoles),
					zap.String("user_email", user.Email))
				return response.Forbidden(c, "Insufficient permissions")
			}

			// Set user role in context for handlers to use
			c.Set("user_role", user.Role)

			logger.Debug("Multi-role RBAC check passed", 
				zap.String("request_id", requestID),
				zap.Int("user_id", userID),
				zap.String("user_role", user.Role),
				zap.Strings("allowed_roles", allowedRoles))

			return next(c)
		}
	}
}

// OwnerOrRoleMiddleware creates middleware that allows access if user owns the resource OR has required role
func OwnerOrRoleMiddleware(userRepo repository.UserRepository, resourceUserIDParam string, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Get user ID from JWT token (set by AuthMiddleware)
			userID, ok := c.Get("user_id").(int)
			if !ok {
				logger.Warn("Owner/Role RBAC check failed: user not authenticated", 
					zap.String("request_id", requestID),
					zap.String("resource_param", resourceUserIDParam),
					zap.Strings("allowed_roles", allowedRoles))
				return response.Unauthorized(c, "User not authenticated")
			}

			// Get resource user ID from URL parameter
			resourceUserIDStr := c.Param(resourceUserIDParam)
			if resourceUserIDStr == "" {
				logger.Warn("Owner/Role RBAC check failed: missing resource parameter", 
					zap.String("request_id", requestID),
					zap.String("resource_param", resourceUserIDParam))
				return response.BadRequest(c, "Invalid request", nil)
			}

			resourceUserID, err := strconv.Atoi(resourceUserIDStr)
			if err != nil {
				logger.Warn("Owner/Role RBAC check failed: invalid resource ID", 
					zap.String("request_id", requestID),
					zap.String("resource_param", resourceUserIDParam),
					zap.String("resource_id", resourceUserIDStr),
					zap.Error(err))
				return response.BadRequest(c, "Invalid resource ID", nil)
			}

			// Check if user owns the resource
			if userID == resourceUserID {
				// User owns the resource, allow access
				logger.Debug("Owner/Role RBAC check passed: resource owner", 
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.Int("resource_user_id", resourceUserID))
				return next(c)
			}

			// User doesn't own the resource, check if they have required role
			user, err := userRepo.GetByID(context.Background(), userID)
			if err != nil {
				if err == sql.ErrNoRows {
					logger.Warn("Owner/Role RBAC check failed: user not found", 
						zap.String("request_id", requestID),
						zap.Int("user_id", userID),
						zap.Strings("allowed_roles", allowedRoles))
					return response.Unauthorized(c, "User not found")
				}
				
				logger.Error("Owner/Role RBAC check failed: database error", 
					zap.Error(err),
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.Strings("allowed_roles", allowedRoles))
				return response.InternalServerError(c, "Internal server error", nil)
			}

			// Check if user has any of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				logger.Warn("Owner/Role RBAC check failed: not owner and insufficient permissions", 
					zap.String("request_id", requestID),
					zap.Int("user_id", userID),
					zap.Int("resource_user_id", resourceUserID),
					zap.String("user_role", user.Role),
					zap.Strings("allowed_roles", allowedRoles),
					zap.String("user_email", user.Email))
				return response.Forbidden(c, "Insufficient permissions")
			}

			// Set user role in context for handlers to use
			c.Set("user_role", user.Role)

			logger.Debug("Owner/Role RBAC check passed: has required role", 
				zap.String("request_id", requestID),
				zap.Int("user_id", userID),
				zap.Int("resource_user_id", resourceUserID),
				zap.String("user_role", user.Role),
				zap.Strings("allowed_roles", allowedRoles))

			return next(c)
		}
	}
}

// AdminMiddleware creates middleware that requires admin role
func AdminMiddleware(userRepo repository.UserRepository) echo.MiddlewareFunc {
	return RoleMiddleware(userRepo, "admin")
}

// ModeratorOrAdminMiddleware creates middleware that requires moderator or admin role
func ModeratorOrAdminMiddleware(userRepo repository.UserRepository) echo.MiddlewareFunc {
	return MultiRoleMiddleware(userRepo, "moderator", "admin")
}

// SelfOrAdminMiddleware creates middleware that allows users to access their own resources or requires admin role
func SelfOrAdminMiddleware(userRepo repository.UserRepository) echo.MiddlewareFunc {
	return OwnerOrRoleMiddleware(userRepo, "id", "admin")
}
