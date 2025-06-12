package helpers

import (
	"testing"

	"github.com/voxpupuli/webhook-go/lib/parsers"
	"gotest.tools/assert"
)

func Test_GetPrefix(t *testing.T) {
	h := Helper{}
	d := parsers.Data{
		RepoName: "testrepo",
		RepoUser: "testuser",
	}

	pfx := "testprefix"

	withPrefix := h.GetPrefix(d, pfx)
	assert.Equal(t, pfx, withPrefix)

	noPrefix := h.GetPrefix(d, "")
	assert.Equal(t, "", noPrefix)

	repoPfx := h.GetPrefix(d, "repo")
	assert.Equal(t, d.RepoName, repoPfx)

	userPfx := h.GetPrefix(d, "user")
	assert.Equal(t, d.RepoUser, userPfx)
}

func Test_GetPrefixFromMapping(t *testing.T) {
	h := Helper{}
	mapping := map[string]string{
		"testrepo":  "testprefix",
		"otherrepo": "otherprefix",
	}

	prefix := h.GetPrefixFromMapping(mapping, "testrepo")
	assert.Equal(t, "testprefix", prefix)

	prefix = h.GetPrefixFromMapping(mapping, "otherrepo")
	assert.Equal(t, "otherprefix", prefix)

	// Test with a repo not in the mapping
	emptyPrefix := h.GetPrefixFromMapping(mapping, "nonexistentrepo")
	assert.Equal(t, "", emptyPrefix)
}
