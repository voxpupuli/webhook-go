package helpers

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func Test_Normalize(t *testing.T) {
	h := Helper{}

	uppercase := true
	lowercase := false
	string := "FooBar"

	upper := h.Normalize(uppercase, string)
	assert.Equal(t, string, upper)

	lower := h.Normalize(lowercase, string)
	assert.Equal(t, strings.ToLower(string), lower)
}
