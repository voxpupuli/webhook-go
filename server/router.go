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
	user := users.Users{
		User:     config.GetConfig().Server.User,
		Password: config.GetConfig().Server.Password,
	}

	health := new(wapi.HealthController)

	router.GET("/health", health.Status)

	api := router.Group("api", gin.BasicAuth(gin.Accounts{user.User: user.Password}))
	{
		v1 := api.Group("v1")
		{
			module := new(wapi.ModuleController)
			v1.POST("/module", module.DeployModule)
			environment := new(wapi.EnvironmentController)
			v1.POST("/environment", environment.DeployEnvironment)
		}
	}

	return router
}
