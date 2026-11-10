package parsers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	bitbucketserver "github.com/go-playground/webhooks/v6/bitbucket-server"
	customerrors "github.com/voxpupuli/webhook-go/lib"
)

// parseBitbucketServer processes a Bitbucket Server webhook, extracting details such as branch, repository, and project info.
// Handles RepositoryReferenceChangedEvent to set branch-related fields.
// When verification is armed, the payload signature is verified by the
// webhook library.
func (d *Data) parseBitbucketServer(c *gin.Context) error {
	secret, err := verificationSecret("bitbucket-server")
	if err != nil {
		return err
	}

	var opts []bitbucketserver.Option
	if secret != "" {
		// Reject a missing or malformed signature up front: the webhook
		// library's HMAC check slices off the "sha256=" prefix and would
		// panic on a shorter header value.
		if !strings.HasPrefix(c.Request.Header.Get("X-Hub-Signature"), "sha256=") {
			return customerrors.ErrWebhookVerificationFailed
		}

		opts = append(opts, bitbucketserver.Options.Secret(secret))
	}

	bh, err := bitbucketserver.New(opts...)
	if err != nil {
		return err
	}

	payload, err := bh.Parse(c.Request, bitbucketserver.RepositoryReferenceChangedEvent)
	if err != nil {
		if secret != "" && (errors.Is(err, bitbucketserver.ErrHMACVerificationFailed) ||
			errors.Is(err, bitbucketserver.ErrMissingHubSignatureHeader)) {
			return customerrors.ErrWebhookVerificationFailed
		}
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
