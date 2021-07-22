package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/voxpupuli/webhook-go/api"
	"github.com/voxpupuli/webhook-go/customerrors"
)

func main() {
	r := api.NewRoutes()

	r.Run(":4000")
}
