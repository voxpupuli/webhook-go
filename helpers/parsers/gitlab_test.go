package parsers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ParseGitlab(t *testing.T) {
	d := Data{}

	header := []Header{
		{
			Name:  "X-Gitlab-Event",
			Value: "Push Hook",
		},
	}

	c, _, err := getGinContext("./json/gitlab.json", header)
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
	}

	assert.NilError(t, err)
	assert.Equal(t, d, d_base)
}
