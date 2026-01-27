package helpers

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
)

const R10kCommand = "r10k"

func (h *Helper) GetR10kCommand() string {
	conf := config.GetConfig().R10k
	commandPath := conf.CommandPath
	if commandPath == "" {
		return R10kCommand
	}
	return commandPath
}

func (h *Helper) GetR10kDeployEnvironmentCommand(env string, noMods bool) []string {
	conf := config.GetConfig().R10k

	cmd := []string{h.GetR10kCommand(), "deploy", "environment"}
	// Append the environment and r10k configuration into the string slice `cmd`
	cmd = append(cmd, env)

	cmd = append(cmd, fmt.Sprintf("--config=%s", h.GetR10kConfig()))

	// Set additional optional r10k options if they are set
	if conf.Verbose {
		cmd = append(cmd, "--verbose")
	}
	if conf.GenerateTypes {
		cmd = append(cmd, "--generate-types")
	}
	// Handle no_mods parameter
	if conf.DeployModules && (noMods) {
		if conf.UseLegacyPuppetfileFlag {
			cmd = append(cmd, "--puppetfile")
		} else {
			cmd = append(cmd, "--modules")
			if conf.EnvironmentIncremental {
				cmd = append(cmd, "--incremental")
			}
		}
	}

	return cmd
}
