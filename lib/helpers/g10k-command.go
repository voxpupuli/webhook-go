package helpers

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
)

const G10kCommand = "g10k"

func (h *Helper) GetG10kCommand() string {
	conf := config.GetConfig().R10k
	commandPath := conf.CommandPath
	if commandPath == "" {
		return G10kCommand
	}
	return commandPath
}

func (h *Helper) GetG10kDeployEnvironmentCommand(env string, branch string) []string {
	conf := config.GetConfig().R10k

	cmd := []string{h.GetG10kCommand()}
	cmd = append(cmd, fmt.Sprintf("-config=%s", h.GetR10kConfig()))

	// Append the environment and g10k configuration into the string slice `cmd`
	if branch == env {
		cmd = append(cmd, fmt.Sprintf("-branch=%s", env))
	} else {
		cmd = append(cmd, fmt.Sprintf("-environment=%s", env))
	}
	// Set additional optional g10k options if they are set
	if conf.Verbose {
		cmd = append(cmd, "-verbose")
	}
	return cmd
}
