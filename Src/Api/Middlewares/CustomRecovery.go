package Middlewares

import (
	"RideMarket-CleanWebApi-GoLang/Api/Helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err any) {
	if err, ok := err.(error); ok {
		httpResponse := Helper.GenerateHttpResponseWithValidationError(nil, false, -500, err.(error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpResponse)
		return
	}
	httpResponse := Helper.GenerateBaseResponseWithAnyError(nil, false, -500, err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, httpResponse)
}
