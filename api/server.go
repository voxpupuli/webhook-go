package api

import (
	"github.com/gin-gonic/gin"
	helpers "github.com/voxpupuli/webhook-go/helpers"
	"net/http"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{}
	router := gin.Default()

	router.GET("/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "running",
		})
	})

	router.POST("/api/r10k/v1/modules", nil)

	server.router = router
	return server
}

func main() {
	//helpers, err := initializeConfig()
	//if err != nil {
	//    panic(err)
	//}

	// authorized := router.Group("/api", gin.BasicAuth(gin.Accounts{
	//     helpers.Authentication.Username: helpers.Authentication.Password,
	// }))
}

func (s *Server) initializeConfig() (*helpers.Config, error) {
	path := helpers.DefaultConfigPath()
	config := &helpers.Config{}

	err := helpers.ValidateConfigPath(path)
	if err != nil {
		return config, err
	}

	config, err = helpers.NewConfig(path)
	if err != nil {
		return config, err
	}

	return config, nil
}
