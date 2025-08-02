package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Initialize() error {
	var config zap.Config
	
	env := os.Getenv("APP_ENV")
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	
	Logger = logger
	return nil
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Sync() {
	Logger.Sync()
}

func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}