package chatops

import (
	"github.com/slack-go/slack"
)

func (c *ChatOps) slack(code int, target string) (*string, *string, error) {
	var sapi *slack.Client
	if c.TestMode {
		sapi = slack.New(c.AuthToken, slack.OptionAPIURL(*c.TestURL+"/"))
	} else {
		sapi = slack.New(c.AuthToken)
	}

	msg := c.formatMessage(code, target)
	attachment := slack.Attachment{
		AuthorName: msg.AuthorName,
		Title:      msg.Title,
		Color:      msg.Color,
		Text:       msg.Text,
	}

	channel, timestamp, err := sapi.PostMessage(c.Channel, slack.MsgOptionUsername(c.User), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return nil, nil, err
	}

	return &channel, &timestamp, nil
}
