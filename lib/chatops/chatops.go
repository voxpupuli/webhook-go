package chatops

import (
	"fmt"
	"strconv"
)

type ChatOps struct {
	Service   string
	Channel   string
	User      string
	AuthToken string
	ServerURI *string
}

type ChatOpsResponse struct {
	Timestamp string
	Channel   string
}

type ChatAttachment struct {
	AuthorName string
	Title      string
	Text       string
	Color      string
}

func (c *ChatOps) PostMessage(code, target string) (*ChatOpsResponse, error) {
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
	default:
		return nil, fmt.Errorf("ChatOps tools `%s` is not supported at this time", c.Service)
	}
	return &resp, nil
}

func (c *ChatOps) formatMessage(code, target string) ChatAttachment {
	var message ChatAttachment

	message.AuthorName = "r10k for Puppet"
	message.Title = fmt.Sprintf("r10k deployment of Puppet environment %s", target)

	if code == "202" {
		message.Text = fmt.Sprintf("Successfully started deployment of %s", target)
		message.Color = "green"
	} else if code == "500" {
		message.Text = fmt.Sprintf("Failed to deploy %s", target)
		message.Color = "red"
	} else {
		message.Text = fmt.Sprintf("Unknown HTTP code: %s", code)
		message.Color = "yellow"
	}

	return message
}
