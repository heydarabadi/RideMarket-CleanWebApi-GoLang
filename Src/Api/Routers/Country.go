package Routers

import (
	"RideMarket-CleanWebApi-GoLang/Api/Handlers"
	"RideMarket-CleanWebApi-GoLang/Config"

	"github.com/gin-gonic/gin"
)

func Country(r *gin.RouterGroup, cfg *Config.Config) {
	h := Handlers.NewCountryHandler(cfg)
	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)

}
