package middleware

import (
	"go-template/internal/logger"
	"go-template/pkg/response"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if requests, exists := rl.requests[identifier]; exists {
		validRequests := make([]time.Time, 0, len(requests))
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[identifier] = validRequests
	}

	// Check if limit exceeded
	if len(rl.requests[identifier]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[identifier] = append(rl.requests[identifier], now)
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for identifier, requests := range rl.requests {
			validRequests := make([]time.Time, 0, len(requests))
			for _, req := range requests {
				if req.After(windowStart) {
					validRequests = append(validRequests, req)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, identifier)
			} else {
				rl.requests[identifier] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimitMiddleware(limiter *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identifier := c.RealIP()

			if !limiter.Allow(identifier) {
				logger.Warn("Rate limit exceeded", zap.String("ip", identifier), zap.String("path", c.Request().URL.Path))
				return response.TooManyRequest(c, "Rate limit exceeded. Please try again later.", nil)
			}

			return next(c)
		}
	}
}
