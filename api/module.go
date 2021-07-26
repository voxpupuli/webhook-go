package api

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

type ModuleController struct{}

func (m ModuleController) DeployModule(c *gin.Context) {
	data := parsers.Data{}

	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		return
	}

	cmd := exec.Command("r10k", "module", "deploy", data.ModuleName)

	out, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deploying module", "error": err})
		c.Abort()
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": out})
}
