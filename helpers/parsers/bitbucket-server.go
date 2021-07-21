package parsers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suhaibmujahid/go-bitbucket-server/bitbucket"
)

func (d *Data) ParseBitbucketServer(c *gin.Context) error {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	event, err := bitbucket.ParseWebHook(bitbucket.WebHookType(c.Request), payload)
	if err != nil {
		return nil
	}

	switch e := event.(type) {
	case *bitbucket.PushEvent:
		d.Branch = d.BsParseBranch(e)
		d.Deleted = d.BitbucketServerDeleted(e)
		d.ModuleName = e.Repository.Name
		d.RepoName = e.Repository.Project.Name + "/" + e.Repository.Name
		d.RepoUser = e.Repository.Project.Name
	default:
		return fmt.Errorf("unknown event type %s", bitbucket.WebHookType(c.Request))
	}

	return nil
}

func (d *Data) BitbucketServerDeleted(c *bitbucket.PushEvent) bool {
	return c.Changes[0].Type == "DELETE"
}

func (d *Data) BsParseBranch(e *bitbucket.PushEvent) string {
	return strings.ReplaceAll(e.Changes[0].Ref.DisplayID, "refs/heads/", "")
}
