package parsers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	api "code.gitea.io/gitea/modules/structs"
	"github.com/gin-gonic/gin"
)

func giteaWebhookType(r *http.Request) string {
	return r.Header.Get("X-Gitea-Event")
}

func (d *Data) parseGitea(c *gin.Context) error {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()
	eventType := giteaWebhookType(c.Request)

	switch eventType {
	// switch e := event.(type) {
	case "push":
		e, err := api.ParsePushHook(payload)
		if err != nil {
			return api.ErrInvalidReceiveHook
		}
		d.Branch = e.Branch()
		// Deletion in Gitea is a different event
		d.Deleted = false
		d.ModuleName = e.Repo.Name
		d.RepoName = e.Repo.FullName
		d.RepoUser = e.Repo.Owner.UserName
		d.Completed = true
		d.Succeed = true
	default:
		return fmt.Errorf("unknown event type %s", giteaWebhookType(c.Request))
	}
	return nil
}
