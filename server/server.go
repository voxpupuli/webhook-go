package server

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/queue"
)

// Init initializes and starts the server on the configured port.
// If queue functionality is enabled in the config, it starts the job queue processing.
// The server can run with or without TLS, depending on the configuration.
func Init() {
	config := config.GetConfig().Server

	if config.Queue.Enabled {
		queue.Work() // Start the job queue if enabled.
	}

	r := NewRouter() // Initialize the router.
	if config.TLS.Enabled {
		// Start the server with TLS (HTTPS) using the provided certificate and key.
		r.RunTLS(":"+fmt.Sprint(config.Port), config.TLS.Certificate, config.TLS.Key)
	} else {
		// Start the server without TLS (HTTP).
		r.Run(":" + fmt.Sprint(config.Port))
	}
}
