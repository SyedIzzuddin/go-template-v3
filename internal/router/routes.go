package router

import (
	"net/http"
	"time"

	"go-template/internal/database"
	"go-template/internal/handler"
	"go-template/internal/middleware"
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
	
	// Protected auth routes
	authProtected := auth.Group("", middleware.AuthMiddleware(jwtManager))
	authProtected.GET("/me", authHandler.GetProfile)

	// Protected user routes
	users := api.Group("/users", middleware.AuthMiddleware(jwtManager))
	users.POST("", userHandler.CreateUser)
	users.GET("", userHandler.GetAllUsers)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)

	// Protected file routes
	files := api.Group("/files", middleware.AuthMiddleware(jwtManager))
	files.POST("/upload", fileHandler.UploadFile)
	files.GET("", fileHandler.GetAllFiles)
	files.GET("/my", fileHandler.GetMyFiles)
	files.GET("/:id", fileHandler.GetFile)
	files.PUT("/:id", fileHandler.UpdateFile)
	files.DELETE("/:id", fileHandler.DeleteFile)
	files.GET("/:id/download", fileHandler.DownloadFile)

	// Static file serving (public)
	e.Static("/uploads", "uploads")
	e.GET("/files/:filename", fileHandler.ServeFile)
}
