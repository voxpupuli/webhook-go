package chatops

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proclaim/mock-slack/server"
)

func Test_PostMessage(t *testing.T) {
	t.Run("ChatOps Message Post", func(t *testing.T) {
		t.Run("Slack", func(t *testing.T) {
			mockServer := server.New()
			c := ChatOps{
				Service:   "slack",
				Channel:   "#general",
				User:      "echo1",
				AuthToken: "12345",
				TestMode:  true,
				TestURL:   &mockServer.Server.URL,
			}

			resp, err := c.PostMessage(202, "main")

			assert.NoError(t, err, "should not error")
			assert.Equal(t, resp.Channel, c.Channel, "channel should be correct")
			assert.NotEmpty(t, resp.Timestamp, "timestamp should not be empty")

			assert.Equal(t, len(mockServer.Received.Attachment), 1)
			assert.Equal(t, mockServer.Received.Attachment[0].Color, "green")
			assert.Equal(t, mockServer.Received.Attachment[0].Text, "Successfully started deployment of main")

			resp, err = c.PostMessage(500, "main")

			assert.NoError(t, err, "should not error")

			assert.Equal(t, len(mockServer.Received.Attachment), 1)
			assert.Equal(t, mockServer.Received.Attachment[0].Color, "red")
			assert.Equal(t, mockServer.Received.Attachment[0].Text, "Failed to deploy main")
		})
		t.Run("RocketChat", func(t *testing.T) {
			c := ChatOps{
				Service:   "rocketchat",
				Channel:   "#general",
				User:      "echo1",
				AuthToken: "12345",
				TestMode:  true,
			}

			_, err := c.PostMessage(202, "main")

			assert.Error(t, err, "A ServerURI must be specified to use RocketChat")

		})
	})
}
