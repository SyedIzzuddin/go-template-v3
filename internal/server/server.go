package server

import (
	"fmt"
	"net/http"
	"time"

	"go-template/internal/config"
	"go-template/internal/database"
	"go-template/internal/handler"
	"go-template/internal/logger"
	"go-template/internal/middleware"
	"go-template/internal/migration"
	"go-template/internal/repository"
	"go-template/internal/router"
	"go-template/internal/service"
	"go-template/pkg/jwt"
	"go-template/pkg/storage"
	"go-template/pkg/validator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *database.DB
}

func NewServer() (*http.Server, error) {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := logger.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run database migrations automatically if enabled
	if cfg.Database.AutoMigrate {
		if err := migration.RunMigrations(db.DB, cfg); err != nil {
			logger.Warn("Failed to run automatic migrations", zap.Error(err))
			// Don't return error here - let the app continue even if migrations fail
			// This allows manual intervention if needed
		}
	} else {
		logger.Info("Auto-migration disabled, skipping database migrations")
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Initialize dependencies
	fileStorage := storage.NewFileStorage()
	validatorInstance := validator.New()
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiresIn,
		cfg.JWT.RefreshExpiresIn,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	fileRepo := repository.NewFileRepository(db.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	fileService := service.NewFileService(fileRepo, fileStorage, cfg)
	authService := service.NewAuthService(userRepo, jwtManager, cfg)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, validatorInstance)
	fileHandler := handler.NewFileHandler(fileService, validatorInstance)
	authHandler := handler.NewAuthHandler(authService, validatorInstance)

	// Initialize middleware
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 100 requests per minute

	// Setup middleware
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.RequestLoggerMiddleware())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.CORSMiddleware())
	e.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Setup routes
	router.SetupRoutes(e, db, userHandler, fileHandler, authHandler, jwtManager)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      e,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Info("Server initialized successfully", zap.Int("port", cfg.Server.Port))

	return httpServer, nil
}
