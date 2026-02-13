package api

import (
	"net/http"
	"slices"
	"strings"

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

	// Get the configuration
	conf := config.GetConfig()

	// Setup chatops connection so we don't have to repeat the process
	conn := helpers.ChatopsSetup()

	// Parse the data from the request and error if the parsing fails
	err := data.ParseData(c)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Error Parsing Webhook", "error": err},
		)
		log.Errorf("error parsing webhook: %s", err)
		c.Abort()
		return
	}

	// Stop the deployment if a pipeline or workflow has failed to run
	// and the DeployOnSuccessOnly setting is set.
	if conf.Server.DeployOnSuccessOnly {
		if err = helpers.GetPipelineStatus(data.Succeed); err != nil {
			c.JSON(
				http.StatusAccepted,
				// GitLab disables webhooks that fail.
				// The webhook is a notification about
				// the status of a pipeline, not a request
				// to deploy an environment.  We accept
				// the notification, and our decision not
				// to deploy on the basis of having received
				// the notification is our own business.
				gin.H{"message": "Declined to deploy the environment", "error": err},
			)
			log.Errorf("didn't deploy environment: %s", err)
			c.Abort()
			return
		}
	}

	// Setup the environment for r10k from the configuration
	if data.Branch == "" {
		branch = conf.R10k.DefaultBranch
	} else {
		branch = data.Branch
	}

	// If branch is listed as a blocked branch, then log it and return.
	if slices.Contains(conf.R10k.BlockedBranches, branch) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"message": "Branch not allowed to be deployed to.", "Branch": branch},
		)
		log.Errorf("branch not permitted for deployment: %s", branch)
		c.Abort()
		return
	}

	prefix := ""
	switch conf.R10k.Prefix {
	case "mapping":
		prefix, err = h.GetPrefixFromMapping(conf.RepoMapping, data.RepoName)

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Error getting Prefix", "error": err},
			)
			log.Errorf("error getting prefix from mapping: %s", err)
			c.Abort()
			return
		}
	default:
		prefix = conf.R10k.Prefix
	}

	env := h.GetEnvironment(branch, prefix, conf.R10k.AllowUppercase)
	noMods := c.Query("no_mods") != "true"

	cmd := []string{}
	if conf.R10k.UseG10kCommands {
		cmd = h.GetG10kDeployEnvironmentCommand(env, branch)
	} else {
		cmd = h.GetR10kDeployEnvironmentCommand(env, noMods)
	}

	log.Printf("Executing Command: %s", strings.Join(cmd, " "))

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
			conn.PostMessage(http.StatusInternalServerError, env, res)
		}
		return
	}

	c.JSON(http.StatusAccepted, res)
	if conf.ChatOps.Enabled {
		conn.PostMessage(http.StatusAccepted, env, res)
	}
}
