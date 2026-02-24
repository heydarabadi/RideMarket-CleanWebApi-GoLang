package Helper

import (
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"net/http"
)

var statusCodeMapping = map[string]int{
	// OTP
	ServiceErrors.OtpExists:   409,
	ServiceErrors.OtpUsed:     409,
	ServiceErrors.OtpNotValid: 400,

	//User
	ServiceErrors.EmailExists:    409,
	ServiceErrors.UserNameExists: 409,
	ServiceErrors.RecordNotFound: 404,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := statusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
