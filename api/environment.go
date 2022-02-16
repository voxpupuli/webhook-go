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

// Environment Controller
type EnvironmentController struct{}

// DeployEnvironment takes in the current Gin context and parses the request
// data into a variable then executes the r10k environment deploy either through
// an orchestrator defined in the Orchestration library or a direct local execution
// of the r10k deploy environment command
func (e EnvironmentController) DeployEnvironment(c *gin.Context) {
	var data parsers.Data
	var h helpers.Helper

	// Set the base r10k command into a slice of strings
	cmd := []string{"r10k", "deploy", "environment"}

	// Get the configuration
	conf := config.GetConfig()

	// Setup chatops connection so we don't have to repeat the process
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

	// Append the environment and r10k configuration into the string slice `cmd`
	cmd = append(cmd, env)
	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k options if they are set
	if conf.R10k.Verbose {
		cmd = append(cmd, "-v")
	}
	if conf.R10k.DeployModules {
		cmd = append(cmd, "-m")
	}
	if conf.R10k.GenerateTypes {
		cmd = append(cmd, "--generate-types")
	}

	// Determine if orchestration is enabled and either pass the cmd string slice to a
	// the orchestrationExec function or localExec function.
	// On an error this will:
	//		* Log the error, orchestration type, and command
	//		* Respond with an HTTP 500 error and return the command result in JSON format
	//		* Abort the request
	//		* Notify ChatOps service if enabled
	//
	// On success this will:
	//		* Respond with an HTTP 202 and the result in JSON format
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
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusAccepted, env)
		}
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
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusAccepted, env)
		}
	}
}
