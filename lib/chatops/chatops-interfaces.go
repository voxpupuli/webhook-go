package chatops

import (
	"github.com/pandatix/gocket-chat/api/chat"
)

type ChatOpsInterface interface {
	PostMessage(code int, target string) (*ChatOpsResponse, error)
	slack(code int, target string) (*string, *string, error)
	rocketChat(code int, target string) (*chat.PostMessageResponse, error)
}
