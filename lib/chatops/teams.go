package chatops

import (
	"fmt"
	"runtime"
	"strings"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	log "github.com/sirupsen/logrus"
)

// Sends a message to a webhook in Microsoft Teams. Returns AdaptiveCard message and error.
func (c *ChatOps) teams(code int, target string, output interface{}) (*adaptivecard.Message, error) {
	var color string
	var status string
	var messageTextBlock adaptivecard.Element
	var detailsBlock adaptivecard.Element
	var details bool

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := *c.ServerURI

	// Get caller function name (either module or environment).
	callerFunction := getCaller()

	// The title for message (first TextBlock element).
	msgTitle := "Deploy " + callerFunction + " " + target

	// Formatted message body.
	msgText := "r10k output:\n\n" + fmt.Sprintf("%s", output)

	// Create blank card.
	card := adaptivecard.NewCard()

	// Determine text color and status text
	if code == 202 {
		if ScanforWarn(output) {
			color = adaptivecard.ColorWarning
			status = "successful with Warnings"
		} else {
			color = adaptivecard.ColorGood
			status = "successful"
		}
	} else {
		color = adaptivecard.ColorAttention
		status = "failed"
	}

	// Create title element.
	headerTextBlock := NewTitleTextBlock(msgTitle, color)

	// Change texts and add details button, depending on type of output (error, QueueItem, any)
	switch fmt.Sprintf("%T", output) {
	case "error":
		messageTextBlock = NewTextBlock(fmt.Sprintf("Error: %s", output), color)
		details = false
	case "*queue.QueueItem":
		messageTextBlock = NewTextBlock(fmt.Sprintf("%s %s added to queue", callerFunction, target), color)
		details = false
	default:
		messageTextBlock = NewTextBlock(fmt.Sprintf("Deployment of %s %s", target, status), color)
		details = true
		detailsBlock = adaptivecard.NewHiddenTextBlock(msgText, true)
		detailsBlock.ID = "detailsBlock"
	}

	// This grouping is used for convenience.
	allTextBlocks := []adaptivecard.Element{
		headerTextBlock,
		messageTextBlock,
	}

	// Add "Details" button to hide/unhide detailed output.
	if details {
		allTextBlocks = append(allTextBlocks, detailsBlock)
		toggleButton := adaptivecard.NewActionToggleVisibility("Details")
		if err := toggleButton.AddTargetElement(nil, detailsBlock); err != nil {
			log.Errorf(
				"failed to add element ID to toggle button: %v",
				err,
			)
			return nil, err
		}

		if err := card.AddAction(true, toggleButton); err != nil {
			log.Errorf(
				"failed to add toggle button action to card: %v",
				err,
			)
			return nil, err
		}
	}

	// Assemble card from all elements.
	if err := card.AddElement(true, allTextBlocks...); err != nil {
		log.Errorf(
			"failed to add text blocks to card: %v",
			err,
		)
		return nil, err
	}

	// Create new Message using Card as input.
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		log.Errorf("failed to create message from card: %v", err)
		return nil, err
	}

	// If not testing, sende the message.
	if !c.TestMode {
		if err := mstClient.Send(webhookUrl, msg); err != nil {
			log.Errorf("failed to send message: %v", err)
			return nil, err
		}
	}

	return msg, err
}

// Scans output for string "WARN". Returns true, if found.
func ScanforWarn(output interface{}) bool {
	asstring := fmt.Sprintf("%s", output)
	return strings.Contains(asstring, "WARN")
}

// Gets caller function. Returns "environment", "module" or empty string.
func getCaller() string {
	var callerFunction string

	pc, _, _, ok := runtime.Caller(3)
	runtimedetails := runtime.FuncForPC(pc)
	if ok && runtimedetails != nil {
		splitStr := strings.Split(runtimedetails.Name(), ".")
		switch splitStr[len(splitStr)-1] {
		case "DeployEnvironment":
			callerFunction = "environment"
		case "DeployModule":
			callerFunction = "module"
		default:
			callerFunction = ""
		}
	}
	return callerFunction
}

// Creates title text block. Returns AdaptiveCard element.
func NewTitleTextBlock(title string, color string) adaptivecard.Element {
	return adaptivecard.Element{
		Type:   adaptivecard.TypeElementTextBlock,
		Wrap:   true,
		Text:   title,
		Style:  adaptivecard.TextBlockStyleHeading,
		Size:   adaptivecard.SizeLarge,
		Weight: adaptivecard.WeightBolder,
		Color:  color,
	}
}

// Creates text block. Returns AdaptiveCard element.
func NewTextBlock(text string, color string) adaptivecard.Element {
	textBlock := adaptivecard.Element{
		Type:  adaptivecard.TypeElementTextBlock,
		Wrap:  true,
		Text:  text,
		Color: color,
	}

	return textBlock
}
