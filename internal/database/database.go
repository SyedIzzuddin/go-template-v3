package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-template/internal/config"
	"go-template/internal/logger"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DB struct {
	*sql.DB
}

func New(cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		cfg.Database.Schema,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Database connection established successfully")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	logger.Info("Closing database connection")
	return db.DB.Close()
}

func (db *DB) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats
	dbStats := db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}