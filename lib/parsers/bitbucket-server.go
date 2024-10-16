package parsers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	bitbucketserver "github.com/go-playground/webhooks/v6/bitbucket-server"
)

// parseBitbucketServer processes a Bitbucket Server webhook, extracting details such as branch, repository, and project info.
// Handles RepositoryReferenceChangedEvent to set branch-related fields.
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

// bitbucketServerDeleted checks if the branch was deleted in the reference change event.
func (d *Data) bitbucketServerDeleted(c bitbucketserver.RepositoryReferenceChangedPayload) bool {
	return c.Changes[0].Type == "DELETE"
}

// bsParseBranch extracts the branch name from the reference change event, removing the ref prefix.
func (d *Data) bsParseBranch(e bitbucketserver.RepositoryReferenceChangedPayload) string {
	return strings.TrimPrefix(e.Changes[0].ReferenceID, prefix)
}
