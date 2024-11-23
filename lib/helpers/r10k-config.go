package helpers

import "github.com/voxpupuli/webhook-go/config"

const ConfigFile = "/etc/puppetlabs/r10k/r10k.yaml"

// GetR10kConfig retrieves the R10k configuration file path.
// If no custom path is set in the configuration, it returns the default path.
func (h *Helper) GetR10kConfig() string {
	conf := config.GetConfig().R10k
	confPath := conf.ConfigPath
	if confPath == "" {
		return ConfigFile
	}
	return confPath
}
