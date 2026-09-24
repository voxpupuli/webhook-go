package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func Test_BranchIgnored(t *testing.T) {
	h := Helper{}

	prefixes := []string{"renovate/", "dependabot/"}

	tests := []struct {
		name     string
		branch   string
		prefixes []string
		expected bool
	}{
		{"matches the first prefix", "renovate/puppet-nfs-4.x", prefixes, true},
		{"matches a later prefix", "dependabot/bundler/rake-13.2.1", prefixes, true},
		{"leaves the default branch alone", "production", prefixes, false},
		{"matches only at the start of the name", "feature/renovate/thing", prefixes, false},
		{"a prefix without a separator still matches", "renovate-lockfile", []string{"renovate"}, true},
		{"no prefixes configured ignores nothing", "renovate/puppet-nfs-4.x", []string{}, false},
		{"nil prefixes ignores nothing", "renovate/puppet-nfs-4.x", nil, false},
		{"an empty prefix does not swallow every branch", "production", []string{""}, false},
		{"an empty branch is never ignored", "", prefixes, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, h.BranchIgnored(tt.branch, tt.prefixes))
		})
	}
}
