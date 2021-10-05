package parsers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
)

func (d *Data) ParseGithub(c *gin.Context) error {
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		return err
	}

	switch e := event.(type) {
	case *github.PushEvent:
		d.Branch = path.Base(*e.Ref)
		d.Deleted = *e.Deleted
		d.ModuleName = *e.Repo.Name
		d.RepoName = *e.Repo.FullName
		d.RepoUser = *e.Repo.Owner.Name
		d.Completed = true
		d.Succeed = true
	case *github.WorkflowRunEvent:
		d.Branch = *e.WorkflowRun.HeadBranch
		d.Deleted = d.githubDeleted(*e.WorkflowRun.HeadSHA)
		d.ModuleName = *e.Repo.Name
		d.RepoName = *e.Repo.FullName
		d.RepoUser = *e.Repo.Owner.Login
		d.Completed = *e.Action == "completed"
		d.Succeed = d.isSucceed(e.WorkflowRun.Conclusion)
	default:
		return fmt.Errorf("unknown event type %s", github.WebHookType(c.Request))
	}
	return nil
}

func (d Data) isSucceed(conclusion *string) bool {
	if conclusion == nil {
		return false
	}

	return *conclusion == "success"
}

func (d *Data) githubDeleted(after string) bool {
	return after == "0000000000000000000000000000000000000000"
}
