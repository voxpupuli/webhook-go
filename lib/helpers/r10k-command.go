package helpers

import "github.com/voxpupuli/webhook-go/config"

const Command = "r10k"

func (h *Helper) GetR10kCommand() string {
	conf := config.GetConfig().R10k
	commandPath := conf.CommandPath
	if commandPath == "" {
		return Command
	}
	return commandPath
}
