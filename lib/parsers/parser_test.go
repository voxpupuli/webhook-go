package parsers

import (
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type Header struct {
	Name  string
	Value string
}

func Test_ParseData(t *testing.T) {
	t.Run("Azure DevOps", func(t *testing.T) {
		t.Run("Successfully Parsed", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Azure-DevOps",
					Value: "git.push",
				},
			}

			c, _, err := getGinContext("./json/azure_devops.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "Fabrikam-Fiber-Git",
				RepoName:   "Fabrikam-Fiber-Git",
				RepoUser:   "278d5cd2-584d-4b63-824a-2ba458937249",
				Completed:  true,
				Succeed:    true,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Failed Parsing", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Azure-DevOps",
					Value: "git.pull",
				},
			}

			c, _, err := getGinContext("./json/azure_devops_fail.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			assert.Error(t, err, "Unknown EventType in webhook payload")
		})
	})
	t.Run("Bitbucket Server", func(t *testing.T) {
		t.Run("Successfully Parsed", func(t *testing.T) {
			d := Data{}

			headers := []Header{
				{
					Name:  "X-Event-Key",
					Value: "repo:refs_changed",
				},
				{
					Name:  "X-Request-Id",
					Value: "abcde12345",
				},
			}

			c, _, err := getGinContext("./json/bitbucket-server.json", headers)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "repository",
				RepoName:   "project/repository",
				RepoUser:   "project",
				Completed:  true,
				Succeed:    true,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Failed to parse", func(t *testing.T) {
			d := Data{}

			headers := []Header{
				{
					Name:  "X-Event-Key",
					Value: "repo:modified",
				},
				{
					Name:  "X-Request-Id",
					Value: "abcde12345",
				},
			}

			c, _, err := getGinContext("./json/bitbucket-server-fail.json", headers)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			assert.Error(t, err, "event not defined to be parsed")
		})
	})
	t.Run("Bitbucket Cloud", func(t *testing.T) {
		t.Run("Successfully Parsed", func(t *testing.T) {
			d := Data{}

			headers := []Header{
				{
					Name:  "X-Event-Key",
					Value: "repo:push",
				},
				{
					Name:  "X-Hook-UUID",
					Value: "aba83f33-f838-4727-aac7-0fc45fac66f7",
				},
			}

			c, _, err := getGinContext("./json/bitbucket-cloud.json", headers)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "hello_app",
				RepoName:   "dhollinger/hello_app",
				RepoUser:   "dhollinger",
				Completed:  true,
				Succeed:    true,
			}

			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Failed to parse", func(t *testing.T) {
			d := Data{}

			headers := []Header{
				{
					Name:  "X-Event-Key",
					Value: "repo:updated",
				},
				{
					Name:  "X-Hook-UUID",
					Value: "aba83f33-f838-4727-aac7-0fc45fac66f7",
				},
			}

			c, _, err := getGinContext("./json/bitbucket-cloud.json", headers)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			assert.Error(t, err, "event not defined to be parsed")
		})
	})
	t.Run("GitLab", func(t *testing.T) {
		t.Run("Successfully Parsed Push", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Gitlab-Event",
					Value: "Push Hook",
				},
			}

			c, _, err := getGinContext("./json/gitlab/push.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "Diaspora",
				RepoName:   "mike/diaspora",
				RepoUser:   "Mike",
				Completed:  true,
				Succeed:    true,
			}

			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Successfully Parsed Pipeline", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Gitlab-Event",
					Value: "Pipeline Hook",
				},
			}

			// Running state
			c, _, err := getGinContext("./json/gitlab/pipeline_running.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  false,
				Succeed:    false,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)

			// Succeed state
			c, _, err = getGinContext("./json/gitlab/pipeline_succeed.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base = Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  true,
				Succeed:    true,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)

			// Failed state
			c, _, err = getGinContext("./json/gitlab/pipeline_failed.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base = Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  true,
				Succeed:    false,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Failed to parse", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Gitlab-Event",
					Value: "Tag Push Hook",
				},
			}

			c, _, err := getGinContext("./json/gitlab/tag_push.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			assert.Error(t, err, "unknown event type Tag Push Hook")
		})
	})
	t.Run("GitHub", func(t *testing.T) {
		t.Run("Successfully Parsed Push", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Github-Event",
					Value: "push",
				},
			}

			c, _, err := getGinContext("./json/github/push.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "simple-tag",
				Deleted:    true,
				ModuleName: "Hello-World",
				RepoName:   "Codertocat/Hello-World",
				RepoUser:   "Codertocat",
				Completed:  true,
				Succeed:    true,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
		t.Run("Successfully Parsed Workflow", func(t *testing.T) {
			d := Data{}

			header := []Header{
				{
					Name:  "X-Github-Event",
					Value: "workflow_run",
				},
			}

			// Running state
			c, _, err := getGinContext("./json/github/workflow_run_running.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base := Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  false,
				Succeed:    false,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)

			// Succeed state
			c, _, err = getGinContext("./json/github/workflow_run_succeed.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base = Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  true,
				Succeed:    true,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)

			// Faled state
			c, _, err = getGinContext("./json/github/workflow_run_failed.json", header)
			if err != nil {
				t.Fatal(err)
			}

			err = d.ParseData(c)

			d_base = Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "webhook-test",
				RepoName:   "meow/webhook-test",
				RepoUser:   "meow",
				Completed:  true,
				Succeed:    false,
			}
			assert.NilError(t, err)
			assert.Equal(t, d, d_base)
		})
	})
}

func getGinContext(filename string, headers []Header) (*gin.Context, *httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	rawjson, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
		Body:   rawjson,
		Method: "POST",
	}

	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}

	c.Request = req

	return c, w, nil
}
