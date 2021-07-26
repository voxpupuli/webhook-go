package server

import (
	"github.com/voxpupuli/webhook-go/config"
)

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(config.GetString("server.port"))
}
