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

// Module Controller
type ModuleController struct{}

// DeployModule takes int the current Gin context and parses the request
// data into a variable then executes the r10k module deploy either through
// a direct local execution of the r10k deploy module command
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		if conf.ChatOps.Enabled {
			conn.PostMessage(http.StatusInternalServerError, "Error Parsing Webhook", err)
		}
		return
	}

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

	module := data.ModuleName
	overrideModule := c.Query("module_name")
	// Restrictions to Puppet module names are: 1) begin with lowercase letter, 2) contain lowercase, digits or underscores
	if overrideModule != "" {
		match, _ := regexp.MatchString("^[a-z][a-z0-9_]*$", overrideModule)
		if !match {
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

	// Append module name and r10k configuration to the cmd string slice
	cmd = append(cmd, module)
	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k flags if they are set
	if conf.R10k.Verbose {
		cmd = append(cmd, "-v")
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
		res, err = queue.AddToQueue("module", data.ModuleName, cmd)
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
			conn.PostMessage(http.StatusInternalServerError, data.ModuleName, err)
		}
		return
	}

	c.JSON(http.StatusAccepted, res)
	if conf.ChatOps.Enabled {
		conn.PostMessage(http.StatusAccepted, data.ModuleName, res)
	}
}
