package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config Config

type Config struct {
	Server struct {
		Protected bool   `mapstructure:"protected"`
		User      string `mapstructure:"user"`
		Password  string `mapstructure:"password"`
		Port      int    `mapstructure:"port,int"`
		TLS       struct {
			Enabled     bool   `mapstructure:"enabled"`
			Certificate string `mapstructure:"certificate"`
			Key         string `mapstructure:"key"`
		} `mapstructure:"tls"`
		Queue struct {
			Enabled           bool `mapstructure:"enabled"`
			MaxConcurrentJobs int  `mapstructure:"max_concurrent_jobs"`
			MaxHistoryItems   int  `mapstructure:"max_history_items"`
		} `mapstructure:"queue"`
	} `mapstructure:"server"`
	ChatOps struct {
		Enabled   bool   `mapstructure:"enabled"`
		Service   string `mapstructure:"service"`
		Channel   string `mapstructure:"channel"`
		User      string `mapstructure:"user"`
		AuthToken string `mapstructure:"auth_token"`
		ServerUri string `mapstructure:"server_uri"`
	} `mapstructure:"chatops"`
	R10k struct {
		CommandPath    string `mapstructure:"command_path"`
		ConfigPath     string `mapstructure:"config_path"`
		DefaultBranch  string `mapstructure:"default_branch"`
		Prefix         string `mapstructure:"prefix"`
		AllowUppercase bool   `mapstructure:"allow_uppercase"`
		Verbose        bool   `mapstructure:"verbose"`
		DeployModules  bool   `mapstructure:"deploy_modules"`
		GenerateTypes  bool   `mapstructure:"generate_types"`
	} `mapstructure:"r10k"`
}

func Init(path *string) {
	var err error
	v := viper.New()

	if path != nil {
		v.SetConfigFile(*path)
	} else {
		v.SetConfigType("yml")
		v.SetConfigName("webhook")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/voxpupuli/webhook/")
		v.AddConfigPath("../config/")
		v.AddConfigPath("config/")
	}
	err = v.ReadInConfig()
	if err != nil {
		log.Fatalf("error on parsing config file: %v", err)
	}

	v = setDefaults(v)

	err = v.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}
}

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("server.port", 4000)
	v.SetDefault("server.protected", false)
	v.SetDefault("server.tls_enabled", false)
	v.SetDefault("server.queue.max_concurrent_jobs", 10)
	v.SetDefault("server.queue.max_history_items", 50)
	v.SetDefault("chatops.enabled", false)
	v.SetDefault("r10k.command_path", "/opt/puppetlabs/puppetserver/bin/r10k")
	v.SetDefault("r10k.config_path", "/etc/puppetlabs/r10k/r10k.yaml")
	v.SetDefault("r10k.default_branch", "master")
	v.SetDefault("r10k.allow_uppercase", false)
	v.SetDefault("r10k.prefix", "")
	v.SetDefault("r10k.verbose", true)
	v.SetDefault("r10k.deploy_modules", true)
	v.SetDefault("r10k.generate_types", true)

	return v
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() Config {
	return config
}
