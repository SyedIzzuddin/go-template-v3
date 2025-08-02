package middleware

import (
	"time"

	"go-template/internal/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// Get request ID from header (set by Echo's request ID middleware)
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			
			// Log request start
			logger.Info("Request started",
				zap.String("request_id", requestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("query", c.Request().URL.RawQuery),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)
			
			// Process request
			err := next(c)
			
			// Calculate duration
			duration := time.Since(start)
			
			// Log request completion
			logger.Info("Request completed",
				zap.String("request_id", requestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("duration", duration),
				zap.Int64("bytes_out", c.Response().Size),
			)
			
			return err
		}
	}
}