package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/voxpupuli/webhook-go/customerrors"
)

type Routes struct {
	router *gin.Engine
}

func NewRoutes() Routes {
	r := Routes{
		router: gin.Default(),
	}

	r.router.GET("/heartbeat", r.heartbeat)

	api := r.router.Group("/api/v1")

	r.addModule(api)

	return r
}

func (r Routes) Run(addr ...string) error {
	return r.router.Run()
}

func (r Routes) heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "running",
	})
}

func handle(f func(c *gin.Context) error) gin.HandlerFunc {
	return func(context *gin.Context) {
		if err := f(context); err != nil {
			if ae, ok := err.(customerrors.AppError); ok {
				context.JSON(ae.StatusCode, gin.H{
					"message": ae.ErrorText,
				})
			} else {
				log.Println(err.Error())
				context.JSON(500, gin.H{
					"message": "Internal server error",
				})
			}
		}

	}
}
