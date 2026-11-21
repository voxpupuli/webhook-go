package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_BranchAllowed(t *testing.T) {
	h := Helper{}

	allowed := []string{"production", "release/*", "feature_??"}

	tests := []struct {
		name     string
		branch   string
		allowed  []string
		expected bool
	}{
		{"exact name", "production", allowed, true},
		{"glob matches one level", "release/1.2", allowed, true},
		{"glob does not cross a slash", "release/1.2/hotfix", allowed, false},
		{"single character wildcard", "feature_ab", allowed, true},
		{"single character wildcard needs exactly two", "feature_abc", allowed, false},
		{"unlisted branch is refused", "development", allowed, false},
		{"a name is not a prefix match", "production2", allowed, false},
		{"empty list allows everything", "development", []string{}, true},
		{"nil list allows everything", "development", nil, true},
		{"empty entry does not allow everything", "development", []string{""}, false},
		{"malformed pattern never matches", "production", []string{"[prod"}, false},
		{"empty branch is left to the caller", "", allowed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, h.BranchAllowed(tt.branch, tt.allowed))
		})
	}
}
