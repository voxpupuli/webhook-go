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

// Determine if orchestration is enabled and either pass the cmd string slice to a
// the orchestrationExec function or localExec function.
//
// This returns an interface of the result of the execution and an error
func execute(cmd []string) (interface{}, error) {
	conf := config.GetConfig()
	var res interface{}
	var err error
	if conf.Orchestration.Enabled {
		res, err = orchestrationExec(cmd)
		if err != nil {
			return res, err
		}
	} else {
		res, err = localExec(cmd)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func orchestrationExec(cmd []string) (interface{}, error) {
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
