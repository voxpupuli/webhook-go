package parsers

import (
	"fmt"
	"path"

	"github.com/gin-gonic/gin"
	bitbucketserver "github.com/go-playground/webhooks/v6/bitbucket-server"
)

func (d *Data) parseBitbucketServer(c *gin.Context) error {
	bh, err := bitbucketserver.New()
	if err != nil {
		return err
	}

	payload, err := bh.Parse(c.Request, bitbucketserver.RepositoryReferenceChangedEvent)
	if err != nil {
		return err
	}

	switch p := payload.(type) {
	case bitbucketserver.RepositoryReferenceChangedPayload:
		d.Branch = d.bsParseBranch(p)
		d.Deleted = d.bitbucketServerDeleted(p)
		d.ModuleName = p.Repository.Name
		d.RepoName = p.Repository.Project.Name + "/" + p.Repository.Name
		d.RepoUser = p.Repository.Project.Name
		d.Completed = true
		d.Succeed = true
	default:
		return fmt.Errorf("unknown event type %s", payload)
	}

	return nil
}

func (d *Data) bitbucketServerDeleted(c bitbucketserver.RepositoryReferenceChangedPayload) bool {
	return c.Changes[0].Type == "DELETE"
}

func (d *Data) bsParseBranch(e bitbucketserver.RepositoryReferenceChangedPayload) string {
	return path.Base(e.Changes[0].ReferenceID)
}
