package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/helpers"
	"github.com/voxpupuli/webhook-go/lib/parsers"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

// Environment Controller
type EnvironmentController struct{}

// DeployEnvironment takes in the current Gin context and parses the request
// data into a variable then executes the r10k environment deploy as direct
// local execution of the r10k deploy environment command
func (e EnvironmentController) DeployEnvironment(c *gin.Context) {
	var data parsers.Data
	var h helpers.Helper
	var branch string

	// Set the base r10k command into a slice of strings
	cmd := []string{h.GetR10kCommand(), "deploy", "environment"}

	// Get the configuration
	conf := config.GetConfig()

	// Setup chatops connection so we don't have to repeat the process
	conn := helpers.ChatopsSetup()

	// Parse the data from the request and error if the parsing fails
	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		log.Errorf("error parsing webhook: %s", err)
		c.Abort()
		return
	}

	// Setup the environment for r10k from the configuration
	if data.Branch == "" {
		branch = conf.R10k.DefaultBranch
	} else {
		branch = data.Branch
	}

	env := h.GetEnvironment(branch, conf.R10k.Prefix, conf.R10k.AllowUppercase)

	// Append the environment and r10k configuration into the string slice `cmd`
	cmd = append(cmd, env)

	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k options if they are set
	if conf.R10k.Verbose {
		cmd = append(cmd, "--verbose")
	}
	if conf.R10k.GenerateTypes {
		cmd = append(cmd, "--generate-types")
	}
	if conf.R10k.DeployModules {
		if conf.R10k.UseLegacyPuppetfileFlag {
			cmd = append(cmd, "--puppetfile")
		} else {
			cmd = append(cmd, "--modules")
		}
	}

	// Pass the command to the execute function and act on the result and any error
	// that is returned
	//
	// On an error this will:
	//		* Log the error and command
	//		* Respond with an HTTP 500 error and return the command result in JSON format
	//		* Abort the request
	//		* Notify ChatOps service if enabled
	//
	// On success this will:
	//		* Respond with an HTTP 202 and the result in JSON format

	var res interface{}
	if conf.Server.Queue.Enabled {
		res, err = queue.AddToQueue("env", env, cmd)
	} else {
		res, err = helpers.Execute(cmd)

		if err != nil {
			log.Errorf("failed to execute local command `%s` with error: `%s` `%s`", cmd, err, res)
		}
	}

	if err != nil {
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

}
