package chatops

import (
	"fmt"
	"net/http"

	gochat "github.com/pandatix/gocket-chat"
	"github.com/pandatix/gocket-chat/api/chat"
)

func (c *ChatOps) rocketChat(code int, target string) (*chat.PostMessageResponse, error) {
	if c.ServerURI == nil {
		return nil, fmt.Errorf("A ServerURI must be specified to use RocketChat")
	}

	client := &http.Client{}

	rc, err := gochat.NewRocketClient(client, *c.ServerURI, c.AuthToken, c.User)
	if err != nil {
		return nil, err
	}

	var attachments []chat.Attachement
	msg := c.formatMessage(code, target)

	attachments = append(attachments, chat.Attachement{
		AuthorName: msg.AuthorName,
		Title:      msg.Title,
		Color:      msg.Color,
		Text:       msg.Text,
	})

	pmp := chat.PostMessageParams{
		Channel:      c.Channel,
		Attachements: &attachments,
	}

	res, err := chat.PostMessage(rc, pmp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
