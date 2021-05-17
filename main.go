package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/voxpupuli/webhook-go/helpers/config"
)

type User struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"` 
}

func main() {
    config, err := initializeConfig()
    if err != nil {
        panic(err)
    }

    router := gin.Default()

    router.GET("/heartbeat", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "running",
        })
    })

    authorized := router.Group("/api", gin.BasicAuth(gin.Accounts{
        config.Authentication.Username: config.Authentication.Password,
    }))

    router.Run()
}

func initializeConfig() (*config.Config, error) {
    path := config.DefaultConfigPath()
    config := &config.Config{}

    err := config.ValidateConfigPath(path)
    if err != nil {
        return config, err
    }

    config, err = config.NewConfig(path)
    if err != nil {
        return config, err
    }

    return config, nil
}

