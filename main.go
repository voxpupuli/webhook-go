package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/voxpupuli/webhook-go/helpers"
)

type User struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"` 
}

func main() {
    router := gin.Default()

    router.GET("/heartbeat", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "running",
        })
    })

    authorized := router.Group("/secrets", gin.BasicAuth(gin.Accounts{
        "dhollinger": "pass",
    }))

    authorized.GET("/test", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "secret": "Now you now the secret!",
        })
    })

    router.Run()
}

func initializeConfig() 
