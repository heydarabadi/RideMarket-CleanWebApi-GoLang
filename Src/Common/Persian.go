package Common

import (
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"regexp"
)

const iranianMobileNumberPattern string = `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9]{7}$`

func IranianMobileNumberValidator(mobileNumber string, logger Log.ILogger) bool {
	res, err := regexp.MatchString(iranianMobileNumberPattern, mobileNumber)
	if err != nil {
		logger.Warning(Log.Validation, Log.MobileValidation, err.Error(), nil)
	}
	return res
}
