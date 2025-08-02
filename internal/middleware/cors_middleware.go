package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://srv683364.hstgr.cloud",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestID,
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
