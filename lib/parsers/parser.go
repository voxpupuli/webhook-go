package parsers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Branch     string
	Deleted    bool
	ModuleName string
	RepoName   string
	RepoUser   string
	Completed  bool
	Succeed    bool
}

// ParseData will takes in a *gin.Context c and parses webhook data into a Data struct.
// This function returns an error if something goes wrong and nil if it completes successfully.
func (d *Data) ParseData(c *gin.Context) error {
	vcs, err := d.ParseHeaders(&c.Request.Header)
	if err != nil {
		return err
	}

	switch vcs {
	case "github":
		err = d.ParseGithub(c)
		if err != nil {
			return err
		}
	case "gitlab":
		err = d.ParseGitlab(c)
		if err != nil {
			return err
		}
	case "bitbucket-cloud":
		err = d.ParseBitbucket(c)
		if err != nil {
			return err
		}
	case "bitbucket-server":
		err = d.ParseBitbucketServer(c)
		if err != nil {
			return err
		}
	case "azuredevops":
		err = d.ParseAzureDevops(c)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported version control systems: %s", vcs)
	}

	return nil
}

// ParseHeaders parses the headers and returns a string containing the VCS tool that is making the request.
// If an unsupported VCS tool makes a request, then an error is returned.
func (d *Data) ParseHeaders(headers *http.Header) (string, error) {
	if headers.Get("X-Github-Event") != "" {
		return "github", nil
	} else if headers.Get("X-Gitlab-Event") != "" {
		return "gitlab", nil
	} else if headers.Get("X-Event-Key") != "" {
		if headers.Get("X-Hook-UUID") != "" {
			return "bitbucket-cloud", nil
		} else if headers.Get("X-Request-Id") != "" {
			return "bitbucket-server", nil
		}
	} else if headers.Get("X-Azure-DevOps") != "" {
		return "azuredevops", nil
	} else {
		return "", errors.New("your Webhook provider is not supported")
	}

	return "", errors.New("couldn't find a valid provider")
}
