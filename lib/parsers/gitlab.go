package parsers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xanzy/go-gitlab"
)

func (d *Data) ParseGitlab(c *gin.Context) error {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	event, err := gitlab.ParseHook(gitlab.HookEventType(c.Request), payload)
	if err != nil {
		return err
	}

	switch e := event.(type) {
	case *gitlab.PushEvent:
		d.Branch = strings.ReplaceAll(e.Ref, "refs/heads/", "")
		d.Deleted = d.GitlabDeleted(e)
		d.ModuleName = e.Project.Name
		d.RepoName = e.Project.PathWithNamespace
		d.RepoUser = e.Project.Namespace
	default:
		return fmt.Errorf("unknown event type %s", gitlab.HookEventType(c.Request))
	}

	return nil
}

func (d *Data) GitlabDeleted(c *gitlab.PushEvent) bool {
	return c.After == "0000000000000000000000000000000000000000"
}
