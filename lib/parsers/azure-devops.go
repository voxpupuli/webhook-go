package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

// parseAzureDevops processes an Azure DevOps webhook, extracting event details such as branch, module name, and repository info.
// It handles the PushEvent type and marks the data as completed and successful upon successful parsing.
func (d *Data) parseAzureDevops(c *gin.Context) error {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	event, err := azuredevops.ParseWebHook(payload)
	if err != nil {
		return err
	}

	switch event.PayloadType {
	case azuredevops.PushEvent:
		parsed, err := d.parseRawResource(event)
		if err != nil {
			return err
		}
		d.Branch = d.parseBranch(parsed)
		d.Deleted = d.azureDevopsDeleted(parsed)
		d.ModuleName = *parsed.Repository.Name
		d.RepoName = *parsed.Repository.Name
		d.RepoUser = *parsed.Repository.ID
		d.Completed = true
		d.Succeed = true
	default:
		return fmt.Errorf("unknown event type %v", event.PayloadType)
	}

	return nil
}

// parseRawResource unmarshals the raw payload of a Git push event from Azure DevOps.
func (d *Data) parseRawResource(e *azuredevops.Event) (payload *azuredevops.GitPush, err error) {
	payload = &azuredevops.GitPush{}

	err = json.Unmarshal(e.RawPayload, &payload)
	if err != nil {
		return nil, err
	}

	e.Resource = payload
	return payload, nil
}

// azureDevopsDeleted checks if the push event represents a branch deletion in Azure DevOps.
func (d *Data) azureDevopsDeleted(e *azuredevops.GitPush) bool {
	return *e.RefUpdates[0].NewObjectID == "0000000000000000000000000000000000000000"
}

// parseBranch extracts the branch name from the push event, removing the ref prefix.
func (d *Data) parseBranch(e *azuredevops.GitPush) string {
	return strings.TrimPrefix(*e.RefUpdates[0].Name, prefix)
}
