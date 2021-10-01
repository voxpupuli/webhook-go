package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/server"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: server -c {path}")
		os.Exit(1)
	}
	confPath := flag.String("c", ".", "")
	flag.Parse()
	config.Init(*confPath)
	server.Init()
}
