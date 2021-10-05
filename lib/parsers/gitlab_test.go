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
}

func Test_ParseGitlabPipelineEvent(t *testing.T) {
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
}
