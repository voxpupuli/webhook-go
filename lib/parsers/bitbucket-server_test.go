package parsers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ParseBitbucketServer(t *testing.T) {
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
}
