package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_GetEnvironment(t *testing.T) {
	h := Helper{}

	branch := "release"
	prefix := "prefix"
	env := "prefix_release"

	envOne := h.GetEnvironment(branch, "", false)
	assert.Equal(t, branch, envOne)

	envTwo := h.GetEnvironment(branch, prefix, false)
	assert.Equal(t, env, envTwo)
}
