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
}

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
	default:
		return fmt.Errorf("unsupported version control systems: %s", vcs)
	}

	return nil
}

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
		return "tfs", nil
	} else {
		return "", errors.New("your Webhook provider is not supported")
	}

	return "", errors.New("couldn't find a valid provider")
}
