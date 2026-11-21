package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_BranchMatches(t *testing.T) {
	h := Helper{}

	patterns := []string{"production", "release-*", "feature_??", "team/*/wip"}

	tests := []struct {
		name     string
		branch   string
		patterns []string
		expected bool
	}{
		{"exact name", "production", patterns, true},
		{"star matches a suffix", "release-1.2", patterns, true},
		{"star matches nothing", "release-", patterns, true},
		{"star also crosses a slash", "team/a/b/wip", patterns, true},
		{"question mark matches one character", "feature_ab", patterns, true},
		{"question mark needs exactly one", "feature_abc", patterns, false},
		{"unlisted branch", "development", patterns, false},
		{"a name is not a prefix match", "production2", patterns, false},
		{"pattern is anchored at the start", "old-release-1", patterns, false},
		{"regexp characters are literal", "release-1.2", []string{"release-1.2"}, true},
		{"a dot is not a wildcard", "release-1x2", []string{"release-1.2"}, false},
		{"empty entry matches nothing", "development", []string{""}, false},
		{"no patterns match nothing", "production", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, h.BranchMatches(tt.branch, tt.patterns))
		})
	}
}

func Test_BranchAllowed(t *testing.T) {
	h := Helper{}

	allowed := []string{"production", "release-*"}

	assert.Equal(t, true, h.BranchAllowed("production", allowed))
	assert.Equal(t, true, h.BranchAllowed("release-1.2", allowed))
	assert.Equal(t, false, h.BranchAllowed("development", allowed))
	assert.Equal(t, true, h.BranchAllowed("development", []string{}), "empty list allows everything")
	assert.Equal(t, true, h.BranchAllowed("development", nil), "nil list allows everything")
	assert.Equal(t, false, h.BranchAllowed("development", []string{""}), "an empty entry does not allow everything")
	assert.Equal(t, true, h.BranchAllowed("", allowed), "an empty branch is left to the caller")
}
