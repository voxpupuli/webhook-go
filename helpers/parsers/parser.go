package parsers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

type Data struct {
	Branch     string
	Deleted    bool
	ModuleName string
	RepoName   string
	RepoUser   *string
}

func ParseData(c *gin.Context) (Data, error) {
	data := Data{}
	vcs, err := parseHeaders(&c.Request.Header)
	if err != nil {
		return Data{}, err
	}

	switch {
	case vcs == "github":
		data, err = parseGithub(c)
		if err != nil {
			return Data{}, err
		}
	}

	return data, nil
}

func parseHeaders(headers *http.Header) (string, error) {
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
	} else if headers.Get("X-Atlassian-Token") != "" {
		return "stash", nil
	} else if headers.Get("X-Azure-DevOps") != "" {
		return "tfs", nil
	} else {
		return "", errors.New("your Webhook provider is not supported")
	}

	return "", errors.New("couldn't find a valid provider")
}

func parseGithub(c *gin.Context) (Data, error) {
	var data Data
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return Data{}, err
	}
	defer c.Request.Body.Close()

	gh, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		return Data{}, err
	}

	switch e := gh.(type) {
	case *github.PushEvent:
		data = Data{
			Branch:     *e.Ref,
			Deleted:    false,
			ModuleName: *e.Repo.Name,
			RepoName:   *e.Repo.Name,
			RepoUser:   e.Repo.Organization,
		}
	default:
		err := fmt.Errorf("unknown event type %s", github.WebHookType(c.Request))
		return Data{}, err
	}

	return data, nil
}
