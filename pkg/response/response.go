package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c echo.Context, message string, err interface{}) error {
	return c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func TooManyRequest(c echo.Context, message string, err interface{}) error {
	return c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

func Conflict(c echo.Context, message string, err interface{}) error {
	return c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func InternalServerError(c echo.Context, message string, err interface{}) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func ValidationError(c echo.Context, message string, validationErrors interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: message,
		Error:   validationErrors,
	})
}
