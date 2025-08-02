package validator

import (
	"go-template/internal/dto"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
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