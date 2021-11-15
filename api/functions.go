package api

import (
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
