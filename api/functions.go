package api

import (
	"os/exec"
	"strings"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/chatops"
	"github.com/voxpupuli/webhook-go/lib/orchestrators"
)

func chatopsSetup() *chatops.ChatOps {
	conf := config.GetConfig().ChatOps
	c := chatops.ChatOps{
		Service:   conf.Service,
		Channel:   conf.Channel,
		User:      conf.User,
		AuthToken: conf.AuthToken,
		ServerURI: &conf.ServerUri,
	}

	return &c
}

func orchestrationExec(cmd []string) (*orchestrators.Result, error) {
	command := "\""
	for i := range cmd {
		command = command + cmd[i] + " "
	}
	command = strings.TrimSuffix(command, " ")
	command = command + "\""

	res, err := orchestrators.Deploy(command)
	if err != nil {
		return res, err
	}

	return res, nil
}

func localExec(cmd []string) ([]byte, error) {
	args := cmd[1:]
	command := exec.Command(cmd[0], args...)

	res, err := command.CombinedOutput()
	if err != nil {
		return res, err
	}

	return res, nil
}
