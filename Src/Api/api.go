package Api

import (
	"RideMarket-CleanWebApi-GoLang/Api/Middlewares"
	"RideMarket-CleanWebApi-GoLang/Api/Routers"
	"RideMarket-CleanWebApi-GoLang/Api/Validations"
	"RideMarket-CleanWebApi-GoLang/Config"
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
	r.Use(gin.Logger(), gin.Recovery())
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
		val.RegisterValidation("mobile", Validations.IranianMobileNumberValidator, true)
		val.RegisterValidation("password", Validations.PasswordValidator, true)
	}
}

func RegisterRoute(r *gin.Engine, cfg *Config.Config) {

	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		Routers.Health(health)

		user := v1.Group("/user/register")
		Routers.SendOtp(user, cfg)

	}

}
