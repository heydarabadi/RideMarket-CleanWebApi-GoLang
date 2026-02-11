package Middlewares

import (
	"RideMarket-CleanWebApi-GoLang/Api/Helper"
	"RideMarket-CleanWebApi-GoLang/Config"
	limiter "RideMarket-CleanWebApi-GoLang/pkg/Limiter"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func OtpLimiter(cfg *Config.Config) gin.HandlerFunc {
	var limiter = limiter.NewIPRateLimiter(rate.Every(cfg.Otp.Limiter*time.Second), 1)
	return func(c *gin.Context) {
		limiter := limiter.GetLimiter(getIP(c.Request.RemoteAddr))
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, Helper.GenerateHttpResponseWithError(nil, false, 403, errors.New("not allowed")))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func getIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}
