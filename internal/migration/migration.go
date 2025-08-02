package migration

import (
	"database/sql"
	"fmt"
	"go-template/internal/config"
	"go-template/internal/logger"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// RunMigrations runs all pending database migrations
func RunMigrations(db *sql.DB, cfg *config.Config) error {
	logger.Info("Running database migrations...")

	// Set the database connection for goose
	goose.SetBaseFS(nil)

	// Run migrations from the migrations directory
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set goose dialect", zap.Error(err))
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Run up migrations
	if err := goose.Up(db, "db/migrations"); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB) error {
	logger.Info("Checking migration status...")

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		logger.Error("Failed to get migration status", zap.Error(err))
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	logger.Info("Current migration version", zap.Int64("version", version))
	return nil
}