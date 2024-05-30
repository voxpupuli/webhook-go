package chatops

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	"github.com/stretchr/testify/assert"
)

func TestTeams(t *testing.T) {
	t.Run("Teams Message Post", func(t *testing.T) {
		serverURI := "https://example.webhook.office.com/webhook/xxx"
		c := ChatOps{
			Service:   "teams",
			TestMode:  true,
			ServerURI: &serverURI,
		}

		t.Run("Good module", func(t *testing.T) {
			// test a successful message
			msg, err := c.teams(202, "good_module", "things went well")

			// Prepare message payload
			msg.Prepare()

			var data adaptivecard.Message
			json.Unmarshal([]byte(msg.PrettyPrint()), &data)

			assert.NoError(t, err, "should not error")

			// assert header text
			// two spaces in the text "Deploy  good_module", because calling function can't be determined
			assert.Equal(t, data.Attachments[0].Content.Body[0].Text, "Deploy  good_module", `header should be "Deploy  good_module"`)
			assert.Equal(t, data.Attachments[0].Content.Body[0].Color, "good", `header color should be "good"`)

			// assert body text
			assert.Equal(t, data.Attachments[0].Content.Body[1].Text, "Deployment of good_module successful", `body text should be "Deployment of good_module successful"`)
			assert.Equal(t, data.Attachments[0].Content.Body[1].Color, "good", `body color should be "good"`)

			// assert details text
			assert.Equal(t, data.Attachments[0].Content.Body[2].Text, "r10k output:\n\nthings went well", `details text should be "r10k output:\n\nthings went well"`)
		})
		t.Run("Failed module", func(t *testing.T) {
			// test a fail message
			msg, err := c.teams(500, "fail_module", "something failed")

			// Prepare message payload
			msg.Prepare()

			var data adaptivecard.Message
			json.Unmarshal([]byte(msg.PrettyPrint()), &data)
			// log.Print(msg.PrettyPrint())
			// log.Print(data.Attachments[0].Content.Body[0].Text)

			assert.NoError(t, err, "should not error")

			// assert header text
			// two spaces in the text "Deploy  fail_module", because calling function can't be determined
			assert.Equal(t, data.Attachments[0].Content.Body[0].Text, "Deploy  fail_module", `header should be "Deploy  fail_module"`)
			assert.Equal(t, data.Attachments[0].Content.Body[0].Color, "attention", `header color should be "attention"`)

			// assert body text
			assert.Equal(t, data.Attachments[0].Content.Body[1].Text, "Deployment of fail_module failed", `body text should be "Deployment of fail_module failed"`)
			assert.Equal(t, data.Attachments[0].Content.Body[1].Color, "attention", `body color should be "attention"`)

			// assert details text
			assert.Equal(t, data.Attachments[0].Content.Body[2].Text, "r10k output:\n\nsomething failed", `details text should be "r10k output:\n\nsomething failed"`)
		})
		t.Run("Module with warnings", func(t *testing.T) {
			// test a message with warning
			msg, err := c.teams(202, "warn_module", "WARN -> there are warnings")

			// Prepare message payload
			msg.Prepare()

			var data adaptivecard.Message
			json.Unmarshal([]byte(msg.PrettyPrint()), &data)
			// log.Print(msg.PrettyPrint())
			// log.Print(data.Attachments[0].Content.Body[0].Text)

			assert.NoError(t, err, "should not error")

			// assert header text
			// two spaces in the text "Deploy  warn_module", because calling function can't be determined
			assert.Equal(t, data.Attachments[0].Content.Body[0].Text, "Deploy  warn_module", `header should be "Deploy  warn_module"`)
			assert.Equal(t, data.Attachments[0].Content.Body[0].Color, "warning", `header color should be "attention"`)

			// assert body text
			assert.Equal(t, data.Attachments[0].Content.Body[1].Text, "Deployment of warn_module successful with Warnings", `body text should be "Deployment of warn_module successful with Warnings"`)
			assert.Equal(t, data.Attachments[0].Content.Body[1].Color, "warning", `body color should be "attention"`)

			// assert details text
			assert.Equal(t, data.Attachments[0].Content.Body[2].Text, "r10k output:\n\nWARN -> there are warnings", `details text should be "r10k output:\n\nWARN -> there are warnings"`)
		})
		t.Run("Error", func(t *testing.T) {
			// test an error state
			msg, err := c.teams(500, "error_module", fmt.Errorf("an error occurred"))

			// Prepare message payload
			msg.Prepare()

			var data adaptivecard.Message
			json.Unmarshal([]byte(msg.PrettyPrint()), &data)

			assert.NoError(t, err, "should not error")

			// assert header text
			// two spaces in the text "Deploy  warn_module", because calling function can't be determined
			assert.Equal(t, data.Attachments[0].Content.Body[0].Text, "Deploy  error_module", `header should be "Deploy  error_module"`)
			assert.Equal(t, data.Attachments[0].Content.Body[0].Color, "attention", `header color should be "attention"`)
			assert.Equal(t, data.Attachments[0].Content.Body[1].Text, "Error: an error occurred", `body should be "Error: an error occurred"`)
		})
	})
}
