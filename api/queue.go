package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

// Queue Controller
type QueueController struct{}

// QueueStatus takes in the current Gin context and show the current queue status
func (q QueueController) QueueStatus(c *gin.Context) {
	c.JSON(http.StatusOK, queue.GetQueueItems())
}
