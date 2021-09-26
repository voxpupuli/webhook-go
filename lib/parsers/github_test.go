package parsers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ParseGithub(t *testing.T) {
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
}

func Test_ParseGithubWorkflowEvent(t *testing.T) {
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
}
