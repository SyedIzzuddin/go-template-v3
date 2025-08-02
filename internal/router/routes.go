package router

import (
	"net/http"
	"time"

	"go-template/internal/database"
	"go-template/internal/handler"
	"go-template/internal/middleware"
	"go-template/internal/repository"
	"go-template/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *database.DB, userHandler *handler.UserHandler, fileHandler *handler.FileHandler, authHandler *handler.AuthHandler, jwtManager *jwt.JWTManager) {
	api := e.Group("/api/v1")

	// Health check (public)
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"database":  db.Health(),
			"timestamp": time.Now().Unix(),
		})
	})

	// Authentication routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.GET("/verify-email", authHandler.VerifyEmail)
	auth.POST("/verify-email", authHandler.VerifyEmail)
	auth.POST("/resend-verification", authHandler.ResendVerificationEmail)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.GET("/reset-password", authHandler.ResetPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)
	
	// Initialize repository for RBAC and email verification middleware
	userRepo := repository.NewUserRepository(db.DB)

	// Protected auth routes
	authProtected := auth.Group("", middleware.AuthMiddleware(jwtManager))
	authProtected.GET("/me", authHandler.GetProfile)

	// Protected user routes with RBAC
	users := api.Group("/users", middleware.AuthMiddleware(jwtManager))
	
	// Admin-only user management
	usersAdmin := users.Group("", middleware.AdminMiddleware(userRepo))
	usersAdmin.POST("", userHandler.CreateUser)                    // Only admin can create users
	usersAdmin.DELETE("/:id", userHandler.DeleteUser)              // Only admin can delete users
	
	// Moderator and admin can view all users
	usersModerator := users.Group("", middleware.ModeratorOrAdminMiddleware(userRepo))
	usersModerator.GET("", userHandler.GetAllUsers)                // Moderator+ can list all users
	
	// Self or admin access for individual user operations
	usersSelf := users.Group("", middleware.SelfOrAdminMiddleware(userRepo))
	usersSelf.GET("/:id", userHandler.GetUser)                     // User can view own profile, admin can view any
	usersSelf.PUT("/:id", userHandler.UpdateUser)                  // User can update own profile, admin can update any

	// Protected file routes with email verification warnings and RBAC
	files := api.Group("/files", 
		middleware.AuthMiddleware(jwtManager),
		middleware.EmailVerificationMiddleware(userRepo))
	
	// All authenticated users can upload and view their own files
	files.POST("/upload", fileHandler.UploadFile)                   // Any authenticated user can upload
	files.GET("/my", fileHandler.GetMyFiles)                       // Any authenticated user can view their own files
	
	// Moderator and admin can view all files
	filesModerator := files.Group("", middleware.ModeratorOrAdminMiddleware(userRepo))
	filesModerator.GET("", fileHandler.GetAllFiles)                // Moderator+ can list all files
	filesModerator.DELETE("/:id", fileHandler.DeleteFile)          // Moderator+ can delete any file
	
	// Individual file operations - all authenticated users can access
	files.GET("/:id", fileHandler.GetFile)                         // Any authenticated user can view file metadata
	files.PUT("/:id", fileHandler.UpdateFile)                      // Any authenticated user can update (handler should check ownership)
	files.GET("/:id/download", fileHandler.DownloadFile)           // Any authenticated user can download (handler should check permissions)

	// Static file serving (public)
	e.Static("/uploads", "uploads")
	e.GET("/files/:filename", fileHandler.ServeFile)
}
