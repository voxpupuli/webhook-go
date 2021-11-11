package chatops

import (
	"github.com/pandatix/gocket-chat/api/chat"
)

type ChatOpsInterface interface {
	PostMessage(code, target string) (*ChatOpsResponse, error)
	slack(code, target string) (*string, *string, error)
	rocketChat(code, target string) (*chat.PostMessageResponse, error)
}
