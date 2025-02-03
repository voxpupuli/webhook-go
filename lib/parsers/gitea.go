package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "code.gitea.io/gitea/modules/structs"
	"github.com/gin-gonic/gin"
)

// giteaWebhookType retrieves the event type from the Gitea webhook request.
func giteaWebhookType(r *http.Request) string {
	return r.Header.Get("X-Gitea-Event")
}

// parseGitea processes a Gitea webhook, extracting branch, repository, and user information.
// Handles "push" and "delete" events to set relevant fields based on the payload.
func (d *Data) parseGitea(c *gin.Context) error {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	eventType := giteaWebhookType(c.Request)

	switch eventType {
	case "push":
		e, err := api.ParsePushHook(payload)
		if err != nil {
			return api.ErrInvalidReceiveHook
		}
		d.Branch = e.Branch()
		d.Deleted = false // Deletion in Gitea is a different event
		d.ModuleName = e.Repo.Name
		d.RepoName = e.Repo.FullName
		d.RepoUser = e.Repo.Owner.UserName
		d.Completed = true
		d.Succeed = true
	case "delete":
		e, err := parseDeleteHook(payload)
		if err != nil {
			return api.ErrInvalidReceiveHook
		}
		d.Branch = e.Ref
		d.Deleted = true
		d.ModuleName = e.Repo.Name
		d.RepoName = e.Repo.FullName
		d.RepoUser = e.Repo.Owner.UserName
		d.Completed = true
		d.Succeed = true
	default:
		return fmt.Errorf("unknown event type %s", eventType)
	}
	return nil
}

// This function parses a Gitea delete event into a struct of type
// api.DeletePayload for use later.
func parseDeleteHook(raw []byte) (*api.DeletePayload, error) {
	hook := new(api.DeletePayload)
	if err := json.Unmarshal(raw, hook); err != nil {
		return nil, err
	}

	switch {
	case hook.Repo == nil:
		return nil, api.ErrInvalidReceiveHook
	case len(hook.Ref) == 0:
		return nil, api.ErrInvalidReceiveHook
	}

	return hook, nil
}
