package api

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/helpers"
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

type ModuleController struct{}

func (m ModuleController) DeployModule(c *gin.Context) {
	data := parsers.Data{}
	h := helpers.Helper{}
	cmd := exec.Command("r10k", "deploy", "module")
	conf := config.GetConfig().R10k
	notify := config.GetConfig().ChatOps.Enabled
	conn := chatopsSetup()

	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		if notify {
			conn.PostMessage(http.StatusInternalServerError, "Error Parsing Webhook")
		}
		return
	}

	cmd.Args = append(cmd.Args, data.ModuleName)
	cmd.Args = append(cmd.Args, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	if conf.Verbose {
		cmd.Args = append(cmd.Args, "-v")
	}

	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("cmd.Run() failed with error %s", string(res))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error executing command", "error": string(res)})
		c.Abort()
		if notify {
			conn.PostMessage(http.StatusInternalServerError, data.ModuleName)
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": string(res)})
	log.Info(fmt.Sprintf("\n%s", string(res)))
	if notify {
		conn.PostMessage(http.StatusAccepted, data.ModuleName)
	}
}
