package parsers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/mcdafydd/go-azuredevops/azuredevops"
)

func (d *Data) ParseAzureDevops(c *gin.Context) error {
	payload, err := ioutil.ReadAll(c.Request.Body)
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
		parsed, err := d.ParseRawResource(event)
		if err != nil {
			return err
		}
		d.Branch = d.ParseBranch(parsed)
		d.Deleted = d.AzureDevopsDeleted(parsed)
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

func (d *Data) ParseRawResource(e *azuredevops.Event) (payload *azuredevops.GitPush, err error) {
	payload = &azuredevops.GitPush{}

	err = json.Unmarshal(e.RawPayload, &payload)
	if err != nil {
		return nil, err
	}

	e.Resource = payload
	return payload, nil
}

func (d *Data) AzureDevopsDeleted(e *azuredevops.GitPush) bool {
	return *e.RefUpdates[0].NewObjectID == "0000000000000000000000000000000000000000"
}

func (d *Data) ParseBranch(e *azuredevops.GitPush) string {
	return path.Base(*e.RefUpdates[0].Name)
}
