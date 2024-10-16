package main

import (
	"github.com/voxpupuli/webhook-go/cmd"
)

// // Main function that starts the application
// func main() {
// 	flag.Usage = func() {
// 		fmt.Println("Usage: server -c {path}")
// 		os.Exit(1)
// 	}
// 	confPath := flag.String("c", ".", "")
// 	flag.Parse()
// 	config.Init(*confPath)
// 	server.Init()
// }

// Main function that starts the application
func main() {
	cmd.Execute()
}
