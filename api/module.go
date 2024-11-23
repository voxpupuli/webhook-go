package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/helpers"
	"github.com/voxpupuli/webhook-go/lib/parsers"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

// ModuleController handles module deployment.
type ModuleController struct{}

// DeployModule handles the deployment of a Puppet module via r10k.
// It parses the incoming webhook data, constructs the r10k command,
// and either queues the deployment or executes it immediately.
func (m ModuleController) DeployModule(c *gin.Context) {
	var data parsers.Data
	var h helpers.Helper

	// Set the base r10k command into a string slice
	cmd := []string{h.GetR10kCommand(), "deploy", "module"}

	// Get the configuration
	conf := config.GetConfig()

	// Setup chatops connection so we don't have to repeat the process
	conn := helpers.ChatopsSetup()

	// Parse the data from the request and error if parsing fails
	err := data.ParseData(c)
	if err != nil {
		// Respond with error if parsing fails, notify ChatOps if enabled.
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusInternalServerError, "Error Parsing Webhook", err)
		}
		return
	}

	// Handle optional branch parameter, fallback to default branch if not provided.
	useBranch := c.Query("branch_only")
	if useBranch != "" {
		branch := ""
		if data.Branch == "" {
			branch = conf.R10k.DefaultBranch
		} else {
			branch = data.Branch
		}
		cmd = append(cmd, "-e")
		cmd = append(cmd, branch)
	}

	// Validate module name with optional override from query parameters.
	module := data.ModuleName
	overrideModule := c.Query("module_name")
	if overrideModule != "" {
		match, _ := regexp.MatchString("^[a-z][a-z0-9_]*$", overrideModule)
		if !match {
			// Invalid module name, respond with error and notify ChatOps.
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid module name"})
			c.Abort()
			err = fmt.Errorf("invalid module name: module name does not match the expected pattern; got: %s, pattern: ^[a-z][a-z0-9_]*$", overrideModule)
			if conf.ChatOps.Enabled {
				conn.PostMessage(http.StatusInternalServerError, "Invalid module name", err)
			}
			return
		}
		module = overrideModule
	}

	// Append module name and r10k configuration to the command string slice.
	cmd = append(cmd, module)
	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k flags if they are enabled.
	if conf.R10k.Verbose {
		cmd = append(cmd, "-v")
	}

	// Execute or queue the command based on server configuration.
	var res interface{}
	if conf.Server.Queue.Enabled {
		res, err = queue.AddToQueue("module", data.ModuleName, cmd)
	} else {
		res, err = helpers.Execute(cmd)

		if err != nil {
			log.Errorf("failed to execute local command `%s` with error: `%s` `%s`", cmd, err, res)
		}
	}

	// Handle error response, notify ChatOps if enabled.
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		c.Abort()
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusInternalServerError, data.ModuleName, res)
		}
		return
	}

	// On success, respond with HTTP 202 and notify ChatOps if enabled.
	c.JSON(http.StatusAccepted, res)
	if conf.ChatOps.Enabled {
		conn.PostMessage(http.StatusAccepted, data.ModuleName, res)
	}
}
