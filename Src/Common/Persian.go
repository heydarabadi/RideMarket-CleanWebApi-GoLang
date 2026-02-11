package Common

import (
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"regexp"
)

const IranianMobilePattern string = `^09(10|11|12|13|14|16|17|18|19|21|22|30|32|33|35|36|37|38|39|90|91|92|93|94|99)[0-9]{7}$`

func IranianMobileNumberValidator(mobileNumber string, logger Log.ILogger) bool {
	res, err := regexp.MatchString(IranianMobilePattern, mobileNumber)
	if err != nil {
		logger.Warning(Log.Validation, Log.MobileValidation, err.Error(), nil)
	}
	return res
}
