package server

import (
	"github.com/voxpupuli/webhook-go/config"
)

// The Init function starts the Server on a specific port
func Init() {
	config := config.GetConfig().Server
	r := NewRouter()
	r.Run(":" + config.Port)
}
