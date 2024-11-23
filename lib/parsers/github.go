package parsers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
)

// parseGithub processes a GitHub webhook, extracting branch, repository, and user information.
// Handles both "push" and "workflow_run" events to set relevant fields based on the payload.
func (d *Data) parseGithub(c *gin.Context) error {
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
		d.Branch = strings.TrimPrefix(*e.Ref, prefix)
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

// isSucceed checks if the conclusion of a workflow run is "success".
func (d Data) isSucceed(conclusion *string) bool {
	if conclusion == nil {
		return false
	}
	return *conclusion == "success"
}

// githubDeleted checks if the specified SHA represents a deleted commit.
func (d *Data) githubDeleted(after string) bool {
	return after == "0000000000000000000000000000000000000000"
}
