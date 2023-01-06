package api

import (
	"os/exec"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/chatops"
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

// This returns an interface of the result of the execution and an error
func execute(cmd []string) (interface{}, error) {
	args := cmd[1:]
	command := exec.Command(cmd[0], args...)

	res, err := command.CombinedOutput()
	if err != nil {
		return string(res), err
	}

	return res, nil
}
