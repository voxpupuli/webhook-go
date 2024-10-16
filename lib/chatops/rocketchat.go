package chatops

import (
	"fmt"
	"net/http"

	gochat "github.com/pandatix/gocket-chat"
	"github.com/pandatix/gocket-chat/api/chat"
)

// rocketChat posts a message to a Rocket.Chat channel using the provided HTTP status code and target environment.
// It returns a PostMessageResponse or an error if the operation fails.
// ServerURI must be provided as part of the ChatOps configuration.
func (c *ChatOps) rocketChat(code int, target string) (*chat.PostMessageResponse, error) {
	// Ensure ServerURI is set before proceeding.
	if c.ServerURI == nil {
		return nil, fmt.Errorf("A ServerURI must be specified to use RocketChat")
	}

	client := &http.Client{}

	// Initialize RocketChat client with the provided ServerURI, AuthToken, and User credentials.
	rc, err := gochat.NewRocketClient(client, *c.ServerURI, c.AuthToken, c.User)
	if err != nil {
		return nil, err
	}

	// Format the message based on the HTTP status code and target environment.
	msg := c.formatMessage(code, target)

	// Prepare attachments for the message.
	var attachments []chat.Attachement
	attachments = append(attachments, chat.Attachement{
		AuthorName: msg.AuthorName,
		Title:      msg.Title,
		Color:      msg.Color,
		Text:       msg.Text,
	})

	// Set the parameters for posting the message to the specified channel.
	pmp := chat.PostMessageParams{
		Channel:      c.Channel,
		Attachements: &attachments,
	}

	// Post the message to RocketChat.
	res, err := chat.PostMessage(rc, pmp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
