package Validations

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationErrors(error error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors

	if errors.As(error, &ve) {
		for _, error := range error.(validator.ValidationErrors) {
			var element ValidationError
			element.Property = error.Field()
			element.Tag = error.Tag()
			element.Message = error.Error()
			element.Value = error.Param()
			validationErrors = append(validationErrors, element)
		}
		return &validationErrors
	}

	return nil
}
