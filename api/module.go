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

// Module Controller
type ModuleController struct{}

// DeployModule takes int the current Gin context and parses the request
// data into a variable then executes the r10k module deploy either through
// an orchestrator defined in the orchestration library or a direct local execution
// of the r10k deploy module command
func (m ModuleController) DeployModule(c *gin.Context) {
	var data parsers.Data
	var h helpers.Helper

	// Set the base r10k command into a string slice
	cmd := []string{"r10k", "deploy", "module"}

	// Get the configuration
	conf := config.GetConfig()

	// Setup chatops connection so we don't have to repeat the process
	conn := chatopsSetup()

	// Parse the data from the request and error if parsing fails
	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusInternalServerError, "Error Parsing Webhook")
		}
		return
	}

	// Append module name and r10k configuration to the cmd string slice
	cmd = append(cmd, data.ModuleName)
	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k flags if they are set
	if conf.R10k.Verbose {
		cmd = append(cmd, "-v")
	}

	// Pass the command to the execute function and act on the result and any error
	// that is returned
	//
	// On an error this will:
	//		* Log the error, orchestration type, and command
	//		* Respond with an HTTP 500 error and return the command result in JSON format
	//		* Abort the request
	//		* Notify ChatOps service if enabled
	//
	// On success this will:
	//		* Respond with an HTTP 202 and the result in JSON format
	res, err := execute(cmd)
	if err != nil {
		if conf.Orchestration.Enabled {
			log.Errorf("orchestrator `%s` failed to execute command `%s` with error: `%s` `%s`", *conf.Orchestration.Type, cmd, err, res)
		} else {
			log.Errorf("failed to execute local command `%s` with error: `%s` `%s`", cmd, err, res)
		}

		c.JSON(http.StatusInternalServerError, res)
		c.Abort()
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusInternalServerError, data.ModuleName)
		}
		return
	}
	c.JSON(http.StatusAccepted, res)
	if conf.ChatOps.Enabled {
		conn.PostMessage(http.StatusAccepted, data.ModuleName)
	}
}
