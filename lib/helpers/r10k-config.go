package helpers

import "github.com/voxpupuli/webhook-go/config"

const ConfigFile = "/etc/puppetlabs/r10k/r10k.yaml"

func (h *Helper) GetR10kConfig() string {
	conf := config.GetConfig().R10k
	confPath := conf.ConfigPath
	if confPath == "" {
		return ConfigFile
	}
	return confPath
}
