package Handlers

import "github.com/gin-gonic/gin"

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, "Working!!")
	return
}
