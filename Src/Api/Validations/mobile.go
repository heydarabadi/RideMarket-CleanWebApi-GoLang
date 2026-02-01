package Validations

import (
	"RideMarket-CleanWebApi-GoLang/Common"

	"github.com/go-playground/validator/v10"
)

func IranianMobileNumberValidator(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if !ok {
		return false
	}
	return Common.IranianMobileNumberValidator(value)
}
