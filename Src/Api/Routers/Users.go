package Routers

import (
	"RideMarket-CleanWebApi-GoLang/Api/Handlers"
	"RideMarket-CleanWebApi-GoLang/Api/Middlewares"
	"RideMarket-CleanWebApi-GoLang/Config"

	"github.com/gin-gonic/gin"
)

func SendOtp(router *gin.RouterGroup,
	cfg *Config.Config) {
	h := Handlers.NewUsersHandler(cfg)
	router.POST("/send-otp", Middlewares.OtpLimiter(cfg), h.SendOtp)
	router.POST("/RegisterByUserName", h.RegisterByUserName)
	router.POST("/LoginByUserName", h.LoginByUserName)
	router.POST("/RegisterLoginByMobileNumber", h.RegisterLoginByMobileNumber)

}
