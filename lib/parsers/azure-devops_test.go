package parsers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ParseAzureDevops(t *testing.T) {
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
}
