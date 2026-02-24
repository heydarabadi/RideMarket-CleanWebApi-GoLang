package Api

import (
	"RideMarket-CleanWebApi-GoLang/Api/Middlewares"
	"RideMarket-CleanWebApi-GoLang/Api/Routers"
	"RideMarket-CleanWebApi-GoLang/Api/Validations"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer(cfg *Config.Config) {

	// Config Validations
	RegisterValidators()

	r := gin.New()
	r.Use(Middlewares.Cors(cfg))
	r.Use(gin.Logger(), gin.CustomRecovery(Middlewares.ErrorHandler))
	r.Use(Middlewares.DefaultStructuredLogger(cfg))

	RegisterRoute(r, cfg)
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		return
	}
}

func RegisterValidators() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", Validations.IranianMobileNumberValidator, true)
		if err != nil {
			return
		}
		err := val.RegisterValidation("password", Validations.PasswordValidator, true)
		if err != nil {
			return
		}
	}
}

func RegisterRoute(r *gin.Engine, cfg *Config.Config) {

	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		Routers.Health(health)

		user := v1.Group("/user/register")
		Routers.SendOtp(user, cfg)

		countries := v1.Group("/countries", Middlewares.Authentication(cfg), Middlewares.Authorization([]string{Constants.AdminRoleName}))
		Routers.Country(countries, cfg)

	}

}
