package helpers

import (
	"os/exec"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/chatops"
)

func ChatopsSetup() *chatops.ChatOps {
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
func Execute(cmd []string) (interface{}, error) {
	var res interface{}
	var err error

	res, err = localExec(cmd)
	return res, err
}

func localExec(cmd []string) (string, error) {
	args := cmd[1:]
	command := exec.Command(cmd[0], args...)

	res, err := command.CombinedOutput()
	if err != nil {
		return string(res), err
	}

	return string(res), nil
}
