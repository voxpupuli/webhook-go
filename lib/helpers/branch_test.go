package helpers

import (
	"testing"

	"github.com/voxpupuli/webhook-go/lib/parsers"
	"gotest.tools/assert"
)

func Test_GetBranch(t *testing.T) {
	h := Helper{}

	d := parsers.Data{
		Deleted: false,
		Branch:  "main",
	}

	branch := h.GetBranch(d, "release")
	assert.Equal(t, d.Branch, branch)

	d2 := parsers.Data{
		Deleted: true,
	}

	branch = h.GetBranch(d2, "release")
	assert.Equal(t, "release", branch)

}
