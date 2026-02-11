package Handlers

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Api/Helper"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *Service.UserService
}

func NewUsersHandler(cfg *Config.Config) *UserHandler {
	service := Service.NewUserService(cfg)
	return &UserHandler{service: service}
}

func (h *UserHandler) SendOtp(c *gin.Context) {
	otp := new(Dtos.GetOtpRequest)
	err := c.ShouldBindJSON(otp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateHttpResponseWithValidationError(nil, false,
			-1, err))
		return
	}
	err = h.service.SendOtp(otp)
	if err != nil {
		c.AbortWithStatusJSON(Helper.TranslateErrorToStatusCode(err),
			Helper.GenerateHttpResponseWithError(nil, false, -1, err))
		return
	}

	// Call Internal Sms Service
	c.JSON(http.StatusCreated, Helper.GenerateHttpResponse(nil, true, 0))
}
