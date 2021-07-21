package parsers

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ParseGithub(t *testing.T) {
	d := Data{}

	header := []Header{
		{
			Name:  "X-Github-Event",
			Value: "push",
		},
	}

	c, _, err := getGinContext("./json/github.json", header)
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
	}
	assert.NilError(t, err)
	assert.Equal(t, d, d_base)
}
