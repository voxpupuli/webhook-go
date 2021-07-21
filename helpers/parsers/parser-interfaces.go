package parsers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/bitbucket"
	"github.com/mcdafydd/go-azuredevops/azuredevops"
	bitbucketserver "github.com/suhaibmujahid/go-bitbucket-server/bitbucket"
	"github.com/xanzy/go-gitlab"
)

type WebhookData interface {
	ParseData(c *gin.Context) error
	ParseHeaders(headers *http.Header) (string, error)
	ParseGithub(c *gin.Context) error
	ParseGitlab(c *gin.Context) error
	GitlabDeleted(c *gitlab.PushEvent) bool
	ParseBitbucket(c *gin.Context) error
	BitbucketDeleted(b bitbucket.RepoPushPayload) bool
	ParseBitbucketServer(c *gin.Context) error
	BitbucketServerDeleted(c *bitbucketserver.PushEvent)
	ParseAzureDevops(c *gin.Context) error
	ParseRawResource(e *azuredevops.Event) (payload *azuredevops.GitPush, err error)
	AzureDevopsDeleted(e *azuredevops.GitPush) bool
	ParseBranch(e *azuredevops.GitPush) string
}
