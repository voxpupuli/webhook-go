package chatops

import (
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	"github.com/pandatix/gocket-chat/api/chat"
)

type ChatOpsInterface interface {
	PostMessage(code int, target string) (*ChatOpsResponse, error)
	slack(code int, target string) (*string, *string, error)
	rocketChat(code int, target string) (*chat.PostMessageResponse, error)
	teams(code int, target string, output interface{}) (*adaptivecard.Message, error)
}
