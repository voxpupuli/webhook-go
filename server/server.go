package server

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
)

// The Init function starts the Server on a specific port
func Init() {
	config := config.GetConfig().Server
	r := NewRouter()
	if config.TLS.Enabled {
		r.RunTLS(":"+fmt.Sprint(config.Port), config.TLS.Certificate, config.TLS.Key)
	} else {
		r.Run(":" + fmt.Sprint(config.Port))
	}
}
