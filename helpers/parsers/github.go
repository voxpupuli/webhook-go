package parsers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
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
		d.Branch = strings.ReplaceAll(*e.Ref, "refs/heads/", "")
		d.Deleted = *e.Deleted
		d.ModuleName = *e.Repo.Name
		d.RepoName = *e.Repo.FullName
		d.RepoUser = *e.Repo.Organization
	default:
		return fmt.Errorf("unknown event type %s", github.WebHookType(c.Request))
	}

	return nil
}
