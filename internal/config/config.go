package config

import (
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Upload   UploadConfig
	Email    EmailConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

type DatabaseConfig struct {
	Host           string
	Port           string
	Database       string
	Username       string
	Password       string
	Schema         string
	AutoMigrate    bool
}

type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type UploadConfig struct {
	MaxFileSize   int64
	AllowedTypes  []string
	UploadPath    string
	BaseURL       string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	BaseURL      string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "go-template"),
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:        getEnv("BLUEPRINT_DB_HOST", "localhost"),
			Port:        getEnv("BLUEPRINT_DB_PORT", "5432"),
			Database:    getEnv("BLUEPRINT_DB_DATABASE", "go_template"),
			Username:    getEnv("BLUEPRINT_DB_USERNAME", "postgres"),
			Password:    getEnv("BLUEPRINT_DB_PASSWORD", "password"),
			Schema:      getEnv("BLUEPRINT_DB_SCHEMA", "public"),
			AutoMigrate: getEnvAsBool("DB_AUTO_MIGRATE", true),
		},
		JWT: JWTConfig{
			AccessSecret:     getEnv("JWT_ACCESS_SECRET", "your-super-secret-access-key-change-this-in-production"),
			RefreshSecret:    getEnv("JWT_REFRESH_SECRET", "your-super-secret-refresh-key-change-this-in-production"),
			AccessExpiresIn:  getEnvAsDuration("JWT_ACCESS_EXPIRES_IN", "30m"),
			RefreshExpiresIn: getEnvAsDuration("JWT_REFRESH_EXPIRES_IN", "168h"), // 7 days
		},
		Server: ServerConfig{
			Port:         getEnvAsInt("PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "10s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
		},
		Upload: UploadConfig{
			MaxFileSize:  getEnvAsInt64("UPLOAD_MAX_FILE_SIZE", 10*1024*1024), // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf", "text/plain"},
			UploadPath:   getEnv("UPLOAD_PATH", "uploads"),
			BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@go-template.com"),
			FromName:     getEnv("SMTP_FROM_NAME", "Go Template"),
			BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour
}