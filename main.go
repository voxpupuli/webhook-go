package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/server"
)

func main() {
	environment := flag.String("e", "develoment", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	server.Init()
}
