package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "running",
	})
}
