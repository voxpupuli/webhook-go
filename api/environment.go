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

type EnvironmentController struct{}

func (e EnvironmentController) DeployEnvironment(c *gin.Context) {
	data := parsers.Data{}
	h := helpers.Helper{}
	cmd := exec.Command("r10k", "deploy", "environment")
	conf := config.GetConfig().R10k
	notify := config.GetConfig().ChatOps.Enabled
	conn := chatopsSetup()

	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		log.Errorf("error parsing webhook: %s", err)
		c.Abort()
		return
	}
	prefix := h.GetPrefix(data, conf.Prefix)
	branch := h.GetBranch(data, conf.DefaultBranch)
	env := h.GetEnvironment(branch, prefix, conf.AllowUppercase)

	cmd.Args = append(cmd.Args, env)
	cmd.Args = append(cmd.Args, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	if conf.Verbose {
		cmd.Args = append(cmd.Args, "-v")
	}
	if conf.DeployModules {
		cmd.Args = append(cmd.Args, "-m")
	}
	if conf.GenerateTypes {
		cmd.Args = append(cmd.Args, "--generate-types")
	}

	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("cmd.Run() failed with error %s", string(res))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error executing command", "error": string(res)})
		c.Abort()
		if notify {
			conn.PostMessage(http.StatusInternalServerError, env)
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": string(res)})
	log.Info(fmt.Sprintf("\n%s", string(res)))
	if notify {
		conn.PostMessage(http.StatusAccepted, env)
	}
}
