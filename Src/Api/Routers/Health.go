package Routers

import (
	"RideMarket-CleanWebApi-GoLang/Api/Handlers"

	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	handler := Handlers.NewHealthHandler()

	r.GET("/", handler.Health)
}
