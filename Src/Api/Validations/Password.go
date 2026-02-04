package Validations

import (
	"RideMarket-CleanWebApi-GoLang/Common"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"

	"github.com/go-playground/validator/v10"
)

var logger = Log.NewLogger(Config.GetConfig())

func PasswordValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		fld.Param()
		return false
	}

	err := Common.ValidatePassword(value)
	if err != nil {
		logger.Warning(Log.Validation, Log.MobileValidation, err.Error(), nil)
		return false
	}
	return true
}
