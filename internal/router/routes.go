package router

import (
	"net/http"
	"time"

	"go-template/internal/database"
	"go-template/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *database.DB, userHandler *handler.UserHandler, fileHandler *handler.FileHandler) {
	api := e.Group("/api/v1")

	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"database":  db.Health(),
			"timestamp": time.Now().Unix(),
		})
	})

	// User routes
	users := api.Group("/users")
	users.POST("", userHandler.CreateUser)
	users.GET("", userHandler.GetAllUsers)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)

	// File routes
	files := api.Group("/files")
	files.POST("/upload", fileHandler.UploadFile)
	files.GET("", fileHandler.GetAllFiles)
	files.GET("/my", fileHandler.GetMyFiles)
	files.GET("/:id", fileHandler.GetFile)
	files.PUT("/:id", fileHandler.UpdateFile)
	files.DELETE("/:id", fileHandler.DeleteFile)
	files.GET("/:id/download", fileHandler.DownloadFile)

	// Static file serving
	e.Static("/uploads", "uploads")
	e.GET("/files/:filename", fileHandler.ServeFile)
}
