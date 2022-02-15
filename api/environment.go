package api

import (
	"fmt"
	"net/http"

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
	cmd := []string{"r10k", "deploy", "environment"}
	conf := config.GetConfig()
	conn := chatopsSetup()

	// Parse the data from the request and error if the parsing fails
	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		log.Errorf("error parsing webhook: %s", err)
		c.Abort()
		return
	}

	// Setup the environment for r10k from the configuration
	env := h.GetEnvironment(conf.R10k.DefaultBranch, conf.R10k.Prefix, conf.R10k.AllowUppercase)

	cmd = append(cmd, env)
	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	if conf.R10k.Verbose {
		cmd = append(cmd, "-v")
	}
	if conf.R10k.DeployModules {
		cmd = append(cmd, "-m")
	}
	if conf.R10k.GenerateTypes {
		cmd = append(cmd, "--generate-types")
	}

	if conf.Orchestration.Enabled {
		res, err := orchestrationExec(cmd)
		if err != nil {
			// TODO: Replace with proper custom error types
			log.Errorf("orchestrator `%s` failed to execute command `%s` with error: `%s`", *conf.Orchestration.Type, cmd, err)
			c.JSON(http.StatusInternalServerError, res)
			c.Abort()
			if conf.ChatOps.Enabled {
				conn.PostMessage(http.StatusInternalServerError, env)
			}
			return
		}
		c.JSON(http.StatusAccepted, res)
	} else {
		res, err := localExec(cmd)
		if err != nil {
			log.Errorf("cmd.Run() failed with error %s", string(res))
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error executing command", "error": string(res)})
			c.Abort()
			if conf.ChatOps.Enabled {
				conn.PostMessage(http.StatusInternalServerError, env)
			}
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": string(res)})
		log.Info(fmt.Sprintf("\n%s", string(res)))
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusAccepted, env)
		}
	}
}
