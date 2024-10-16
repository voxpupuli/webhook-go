package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voxpupuli/webhook-go/config"
)

// cfgFile is the path to the configuration file, given by the user as a flag
var cfgFile string

// version is the current version of the application
var version = "0.0.0"

// rootCmd is the root command for the application
// It is used to set up the application, and is the entry point for the Cobra CLI
var rootCmd = &cobra.Command{
	Use:     "webhook-go",
	Version: version,
	Short:   "API Server for providing r10k/g10k as a web service",
	Long: `Provides an API service that parses git-based webhook
	requests, executing r10k deployments based on the payload and
	API endpoint.`,
}

// Execute is the main entry point for the application, called from main.go, and is used to execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init is called when the package loads, and is used to set up the root command, and the configuration file flag
func init() {
	cobra.OnInitialize(initConfig) //  tells Cobra to call the initConfig function before executing any command.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./webhook.yml)") //  adds a flag to the root command that allows the user to specify a configuration file
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		config.Init(&cfgFile) // Expecting a path to a configuration file
	} else {
		config.Init(nil) // No path given, use defaults
	}
}
