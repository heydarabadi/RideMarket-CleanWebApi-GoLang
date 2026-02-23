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

func GenerateBaseHttpResponse(result any, success bool, resultCode int) *BaseHttpResponse {
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

func GenerateBaseResponseWithAnyError(result any, success bool, resultCode int, err any) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:    success,
		ResultCode: resultCode,
		Error:      err,
	}
}

func GenerateBaseResponseWithValidationError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	return &BaseHttpResponse{Result: result,
		Success:          success,
		ResultCode:       resultCode,
		ValidationErrors: Validations.GetValidationErrors(err),
	}
}
