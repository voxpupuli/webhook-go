package parsers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_Bitbucket(t *testing.T) {
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
}
