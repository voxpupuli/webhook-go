/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voxpupuli/webhook-go/server"
)

// serverCmd starts the Webhook-go server, allowing it to process webhook requests.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Webhook-go server",
	Run:   startServer,
}

// init adds serverCmd to the root command.
func init() {
	rootCmd.AddCommand(serverCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// startServer initializes and starts the server.
func startServer(cmd *cobra.Command, args []string) {
	server.Init()
}
