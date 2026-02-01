package Helper

import (
	"RideMarket-CleanWebApi-GoLang/Api/Validations"
)

type BaseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	ResultCode       int                            `json:"resultCode"`
	Error            any                            `json:"error"`
	ValidationErrors *[]Validations.ValidationError `json:"ValidationErrors"`
}

func GenerateHttpResponse(result any, success bool, resultCode int) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
	}
}
func GenerateHttpResponseWithError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
		Error:      err.Error(),
	}
}

func GenerateHttpResponseWithValidationError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:          success,
		ResultCode:       resultCode,
		Error:            err.Error(),
		ValidationErrors: Validations.GetValidationErrors(err),
	}
}
