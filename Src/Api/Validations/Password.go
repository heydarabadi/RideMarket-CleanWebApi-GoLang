package Validations

import (
	"RideMarket-CleanWebApi-GoLang/Common"
	"log"

	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		fld.Param()
		return false
	}

	err := Common.ValidatePassword(value)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	return true
}
