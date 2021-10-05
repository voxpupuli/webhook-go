package parsers

import (
	"fmt"
	"io/ioutil"
	"path"

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
		d.Branch = path.Base(e.Ref)
		d.Deleted = d.GitlabDeleted(e.After)
		d.ModuleName = e.Project.Name
		d.RepoName = e.Project.PathWithNamespace
		d.RepoUser = e.Project.Namespace
		d.Completed = true
		d.Succeed = true
	case *gitlab.PipelineEvent:
		d.Branch = e.ObjectAttributes.Ref
		d.Deleted = d.GitlabDeleted(e.ObjectAttributes.SHA)
		d.ModuleName = e.Project.Name
		d.RepoName = e.Project.PathWithNamespace
		d.RepoUser = e.Project.Namespace
		d.Completed = e.ObjectAttributes.Status != "running"
		d.Succeed = e.ObjectAttributes.Status == "success"
	default:
		return fmt.Errorf("unknown event type %s", gitlab.HookEventType(c.Request))
	}

	return nil
}

func (d *Data) GitlabDeleted(after string) bool {
	return after == "0000000000000000000000000000000000000000"
}
