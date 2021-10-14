package parsers

import (
	"github.com/gin-gonic/gin"
)

type WebhookData interface {
	ParseData(c *gin.Context) error
}
