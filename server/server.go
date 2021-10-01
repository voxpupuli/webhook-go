package server

import (
	"github.com/voxpupuli/webhook-go/config"
)

func Init() {
	config := config.GetConfig().Server
	r := NewRouter()
	r.Run(":" + config.Port)
}
