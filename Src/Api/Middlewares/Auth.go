package Middlewares

import (
	"RideMarket-CleanWebApi-GoLang/Api/Helper"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Constants"
	"RideMarket-CleanWebApi-GoLang/Service"
	"RideMarket-CleanWebApi-GoLang/pkg/ServiceErrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authentication(cfg *Config.Config) gin.HandlerFunc {
	var tokenService = Service.NewTokenService(cfg)
	return func(c *gin.Context) {
		var err error
		claimMaps := map[string]interface{}{}
		auth := c.GetHeader(Constants.AuthorizationHeaderKey)

		token := strings.Split(auth, " ")[1]
		if auth == "" {
			err = &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.TokenRequired}
		} else {
			claimMaps, err = tokenService.GetClaims(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.TokenExpired}
				default:
					err = &ServiceErrors.ServiceError{EndUserMessage: ServiceErrors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				Helper.GenerateHttpResponseWithError(nil, false, -2, err))
			return
		}

		c.Set("userid", claimMaps["userid"])
		c.Set("fullname", claimMaps["fullname"])
		c.Set("username", claimMaps["username"])
		c.Set("email", claimMaps["email"])
		c.Set("role", claimMaps["role"])
		c.Set("exp", claimMaps["exp"])
		c.Set("mobilenumber", claimMaps["mobilenumber"])

		c.Next()
	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if len(context.Keys) == 0 {
			context.AbortWithStatusJSON(http.StatusForbidden, Helper.GenerateBaseHttpResponse(nil, false, -3))
			return
		}
		rolesVal := context.Keys[Constants.RolesKey]
		if rolesVal == nil {
			context.AbortWithStatusJSON(http.StatusForbidden, Helper.GenerateBaseHttpResponse(nil, false, -3))
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				context.Next()
			}
		}
		context.AbortWithStatusJSON(http.StatusForbidden, Helper.GenerateBaseHttpResponse(nil, false, -301))
	}
}
