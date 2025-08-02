package validator

import (
	"go-template/internal/dto"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := validator.New()
	
	// Register custom password validation
	v.RegisterValidation("password", validatePassword)
	
	return &Validator{
		validator: v,
	}
}

// validatePassword checks if password meets security requirements
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}
	
	// Must contain at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}
	
	// Must contain at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}
	
	// Must contain at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return false
	}
	
	// Must contain at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false
	}
	
	return true
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *Validator) ValidateStruct(i interface{}) []dto.ValidationError {
	err := v.validator.Struct(i)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			validationErrors = errors
		}
		return dto.FormatValidationErrors(validationErrors)
	}
	return nil
}