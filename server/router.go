package server

import (
	"github.com/gin-gonic/gin"
	wapi "github.com/voxpupuli/webhook-go/api"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/users"
)

// The NewRouter function sets up the main web routes
// for the Gin API server.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	var apiHandlerFuncs []gin.HandlerFunc
	if config.GetConfig().Server.Protected {
		user := users.Users{
			User:     config.GetConfig().Server.User,
			Password: config.GetConfig().Server.Password,
		}

		apiHandlerFuncs = append(apiHandlerFuncs, gin.BasicAuth(gin.Accounts{user.User: user.Password}))
	}

	health := new(wapi.HealthController)

	router.GET("/health", health.Status)

	api := router.Group("api", apiHandlerFuncs...)
	{
		v1 := api.Group("v1")
		{
			r10k := v1.Group("r10k")
			{
				module := new(wapi.ModuleController)
				r10k.POST("/module", module.DeployModule)
				environment := new(wapi.EnvironmentController)
				r10k.POST("/environment", environment.DeployEnvironment)
			}
		}
	}

	return router
}
