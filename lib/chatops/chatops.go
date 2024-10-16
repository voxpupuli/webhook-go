package chatops

import (
	"fmt"
	"strconv"
)

// ChatOps defines the configuration for interacting with various chat services.
type ChatOps struct {
	Service   string  // Chat service (e.g., "slack", "rocketchat", "teams").
	Channel   string  // Target channel or room.
	User      string  // User initiating the action.
	AuthToken string  // Authentication token for the chat service.
	ServerURI *string // Optional server URI for self-hosted services.
	TestMode  bool    // Indicates if the operation is in test mode.
	TestURL   *string // URL for testing purposes, if applicable.
}

// ChatOpsResponse captures the response details from a chat service after a message is posted.
type ChatOpsResponse struct {
	Timestamp string
	Channel   string
}

// ChatAttachment represents the structure of a message attachment in chat services like Slack.
type ChatAttachment struct {
	AuthorName string
	Title      string
	Text       string
	Color      string // Color to indicate status (e.g., success, failure).
}

// PostMessage sends a formatted message to the configured chat service based on the HTTP status code
// and target environment. It returns a ChatOpsResponse or an error if posting fails.
// Supports Slack, Rocket.Chat, and Microsoft Teams.
func (c *ChatOps) PostMessage(code int, target string, output interface{}) (*ChatOpsResponse, error) {
	var resp ChatOpsResponse

	switch c.Service {
	case "slack":
		ch, ts, err := c.slack(code, target)
		if err != nil {
			return nil, err
		}
		resp.Channel = *ch
		resp.Timestamp = *ts
	case "rocketchat":
		res, err := c.rocketChat(code, target)
		if err != nil {
			return nil, err
		}
		resp.Channel = res.Channel
		resp.Timestamp = strconv.FormatInt(res.Ts, 10)
	case "teams":
		_, err := c.teams(code, target, output)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("ChatOps tools `%s` is not supported at this time", c.Service)
	}
	return &resp, nil
}

// formatMessage generates a ChatAttachment based on the HTTP status code and target environment.
// The message is used to notify the result of a Puppet environment deployment.
func (c *ChatOps) formatMessage(code int, target string) ChatAttachment {
	var message ChatAttachment

	message.AuthorName = "r10k for Puppet"
	message.Title = fmt.Sprintf("r10k deployment of Puppet environment %s", target)

	if code == 202 {
		message.Text = fmt.Sprintf("Successfully started deployment of %s", target)
		message.Color = "green"
	} else if code == 500 {
		message.Text = fmt.Sprintf("Failed to deploy %s", target)
		message.Color = "red"
	} else {
		message.Text = fmt.Sprintf("Unknown HTTP code: %d", code)
		message.Color = "yellow"
	}

	return message
}
