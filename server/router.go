package server

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/voxpupuli/webhook-go/api"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/users"
)

// NewRouter sets up the main routes for the Gin API server.
// It includes logging, recovery middleware, health checks, and API routes for version 1.
// If server protection is enabled, it adds BasicAuth middleware for the API.
func NewRouter() *gin.Engine {
	router := gin.New()        // Create a new Gin router.
	router.Use(gin.Logger())   // Add logging middleware.
	router.Use(gin.Recovery()) // Add recovery middleware to handle panics.

	var apiHandlerFuncs []gin.HandlerFunc
	if config.GetConfig().Server.Protected {
		// If the server is protected, set up BasicAuth using the configured username and password.
		user := users.Users{
			User:     config.GetConfig().Server.User,
			Password: config.GetConfig().Server.Password,
		}
		apiHandlerFuncs = append(apiHandlerFuncs, gin.BasicAuth(gin.Accounts{user.User: user.Password}))
	}

	health := new(wapi.HealthController)
	router.GET("/health", health.Status) // Route for health status check.

	// Group API routes under /api with optional BasicAuth.
	api := router.Group("api", apiHandlerFuncs...)
	{
		v1 := api.Group("v1")
		{
			r10k := v1.Group("r10k") // Group r10k-related routes.
			{
				module := new(wapi.ModuleController)
				r10k.POST("/module", module.DeployModule) // Route for deploying a module.
				environment := new(wapi.EnvironmentController)
				r10k.POST("/environment", environment.DeployEnvironment) // Route for deploying an environment.
			}

			queue := v1.Group("queue") // Group queue-related routes.
			{
				q := new(wapi.QueueController)
				queue.GET("", q.QueueStatus) // Route to check the queue status.
			}
		}
	}

	return router
}
