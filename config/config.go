package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config Config

type Config struct {
	Server struct {
		Protected bool   `yaml:"protected"`
		User      string `yaml:"user"`
		Passsword string `yaml:"password"`
		Port      string `yaml:"port,int"`
	} `yaml:"server"`
	ChatOps struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"chatops"`
	R10k struct {
		ConfigPath     string `yaml:"config_path"`
		DefaultBranch  string `yaml:"default_branch"`
		Prefix         string `yaml:"prefix"`
		AllowUppercase bool   `yaml:"allow_uppercase"`
		Verbose        bool   `yaml:"verbose"`
	} `yaml:"r10k"`
	Pipeline struct {
		Enabled       bool `yaml:"enabled"`
		DeployOnError bool `yaml:"deploy_on_error"`
	} `yaml:"pipeline"`
}

func Init(path string) {
	var err error
	var config Config

	v := viper.New()
	v.SetConfigType("yml")
	v.SetConfigName("webhook")
	v.AddConfigPath(path)
	v.AddConfigPath("/etc/voxpupupli/webhook/")
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
