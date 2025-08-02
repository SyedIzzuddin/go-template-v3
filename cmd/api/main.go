package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-template/internal/logger"
	"go-template/internal/server"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.Initialize(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting application")

	// Initialize server
	srv, err := server.NewServer()
	if err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info("Server closed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
