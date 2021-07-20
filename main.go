package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    helpers "github.com/voxpupuli/webhook-go/helpers"
)

func main() {
    //helpers, err := initializeConfig()
    //if err != nil {
    //    panic(err)
    //}

    router := gin.Default()

    router.GET("/heartbeat", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "running",
        })
    })

    // authorized := router.Group("/api", gin.BasicAuth(gin.Accounts{
    //     helpers.Authentication.Username: helpers.Authentication.Password,
    // }))

    router.GET("/headers", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": c.Request.Header,
        })
    })

    router.Run()
}

func initializeConfig() (*helpers.Config, error) {
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

