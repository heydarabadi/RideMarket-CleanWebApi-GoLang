package Api

import (
	"RideMarket-CleanWebApi-GoLang/Api/Routers"
	"RideMarket-CleanWebApi-GoLang/Config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	cfg := Config.GetConfig()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		Routers.Health(health)
	}

	err := r.Run(fmt.Sprint(cfg.Server.Port))
	if err != nil {
		return
	}
}
