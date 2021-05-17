package helpers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Authentication struct {
        Username string `yaml:"user"`
        Password string `yaml:"pass"`
    } `yaml:"authentication"`
    LogLevel string `yaml:"loglevel"`
}

func NewConfig(configPath string) (*Config, error) {
    config := &Config{}

    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    d := yaml.NewDecoder(file)

    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}

func ValidateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }
    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }
    return nil
}

func DefaultConfigPath() string {
    path := os.Getenv("PWH_CONFIG")
    if path != "" {
        return "/etc/voxpupuli/webhook/config.yml"
    }

    return path
}
