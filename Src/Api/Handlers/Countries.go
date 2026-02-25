package Handlers

import (
	"RideMarket-CleanWebApi-GoLang/Api/Dtos"
	"RideMarket-CleanWebApi-GoLang/Api/Helper"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	service *Service.CountryService
}

func NewCountryHandler(cfg *Config.Config) *CountryHandler {
	return &CountryHandler{service: Service.NewCountryService(cfg)}
}

// Create
func (h *CountryHandler) Create(c *gin.Context) {
	req := Dtos.CreateUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateHttpResponseWithValidationError(nil, false, -7, err))
		return
	}

	res, err := h.service.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(Helper.TranslateErrorToStatusCode(err), Helper.GenerateHttpResponseWithError(nil, false, 111, err))
		return
	}
	c.JSON(http.StatusCreated, Helper.GenerateBaseHttpResponse(res, true, 0))
}

//Update

func (h *CountryHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	req := Dtos.CreateUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateHttpResponseWithValidationError(nil, false, -7, err))
		return
	}

	res, err := h.service.Update(c, id, &req)
	if err != nil {
		c.AbortWithStatusJSON(Helper.TranslateErrorToStatusCode(err), Helper.GenerateHttpResponseWithError(nil, false, 0, err))
		return
	}
	c.JSON(http.StatusOK, Helper.GenerateBaseHttpResponse(res, true, 0))

}

//Delete

func (h *CountryHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateBaseHttpResponse(nil, false, -7))
		return
	}
	err := h.service.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(Helper.TranslateErrorToStatusCode(err), Helper.GenerateHttpResponseWithError(nil, false, 0, err))
		return
	}
	c.JSON(http.StatusOK, Helper.GenerateBaseHttpResponse(nil, true, 0))
}

// GetById
func (h *CountryHandler) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateBaseHttpResponse(nil, false, -7))
		return
	}
	var response *Dtos.CountryResponse
	response, err := h.service.GetById(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Helper.GenerateBaseHttpResponse(nil, false, 0))
		return
	}
	c.JSON(http.StatusOK, Helper.GenerateBaseHttpResponse(response, true, 0))

}

// GetByFilter
func (h *CountryHandler) GetByFilter(c *gin.Context) {
	req := Dtos.PaginationInputWithFilter{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Helper.GenerateHttpResponseWithValidationError(nil, false, -7, err))
		return
	}

	res, err := h.service.GetByFilter(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(Helper.TranslateErrorToStatusCode(err), Helper.GenerateHttpResponseWithError(nil, false, 111, err))
		return
	}
	c.JSON(http.StatusOK, Helper.GenerateBaseHttpResponse(res, true, 0))
}
