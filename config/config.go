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
		Port      string `mapstructure:"port,int"`
		TLS       struct {
			Enabled     bool   `mapstructure:"enabled"`
			Certificate string `mapstructure:"certificate"`
			Key         string `mapstructure:"key"`
		} `mapstructure:"tls"`
	} `mapstructure:"server"`
	ChatOps struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"chatops"`
	R10k struct {
		ConfigPath     string `mapstructure:"config_path"`
		DefaultBranch  string `mapstructure:"default_branch"`
		Prefix         string `mapstructure:"prefix"`
		AllowUppercase bool   `mapstructure:"allow_uppercase"`
		Verbose        bool   `mapstructure:"verbose"`
	} `mapstructure:"r10k"`
	Pipeline struct {
		Enabled       bool `mapstructure:"enabled"`
		DeployOnError bool `mapstructure:"deploy_on_error"`
	} `mapstructure:"pipeline"`
}

func Init(path string) {
	var err error

	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName("webhook")
	v.AddConfigPath(path)
	v.AddConfigPath("/etc/voxpupuli/webhook/")
	v.AddConfigPath("../config/")
	v.AddConfigPath("config/")
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
	v.SetDefault("chatops.enabled", false)
	v.SetDefault("r10k.config_path", "/etc/puppetlabs/r10k/r10k.yaml")
	v.SetDefault("r10k.default_branch", "master")
	v.SetDefault("r10k.allow_uppercase", false)
	v.SetDefault("r10k.prefix", "")
	v.SetDefault("r10k.verbose", true)
	v.SetDefault("pipeline.enabled", false)
	v.SetDefault("pipeline.deploy_on_error", false)

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
