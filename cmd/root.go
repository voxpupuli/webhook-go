package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voxpupuli/webhook-go/config"
)

var cfgFile string

var version = "0.0.0"

var rootCmd = &cobra.Command{
	Use:     "webhook-go",
	Version: version,
	Short:   "API Server for providing r10k/g10k as a web service",
	Long: `Provides an API service that parses git-based webhook
	requests, executing r10k deployments based on the payload and
	API endpoint.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./webhook.yml)")
}

func initConfig() {
	if cfgFile != "" {
		config.Init(&cfgFile)
	} else {
		config.Init(nil)
	}
}
