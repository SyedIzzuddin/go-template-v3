package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func FormatValidationErrors(errs validator.ValidationErrors) []ValidationError {
	var validationErrors []ValidationError
	
	for _, err := range errs {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Param(),
			Message: getValidationMessage(err),
		})
	}
	
	return validationErrors
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "oneof":
		return err.Field() + " must be one of: " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}