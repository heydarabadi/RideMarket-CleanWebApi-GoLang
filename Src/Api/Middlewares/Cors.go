package Middlewares

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(cfg *Config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.Cors.AllowedOrigins)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Max-Age", "21600")
		c.Set("content-type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
