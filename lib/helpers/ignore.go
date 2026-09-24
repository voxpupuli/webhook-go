package helpers

import (
	"strings"
)

// BranchIgnored reports whether branch starts with any of the given prefixes.
//
// This mirrors r10k's own `ignore_branch_prefixes` setting in r10k.yaml. A branch
// r10k is configured to ignore is not an environment at all, so asking r10k to
// deploy it fails with:
//
//	ERROR -> Environment(s) '<name>' cannot be found in any source and will not be deployed.
//
// Callers should skip the deploy and answer with a success status, because the
// sender did nothing wrong: there is simply nothing to deploy.
func (h *Helper) BranchIgnored(branch string, prefixes []string) bool {
	if branch == "" {
		return false
	}
	for _, p := range prefixes {
		if p == "" {
			continue
		}
		if strings.HasPrefix(branch, p) {
			return true
		}
	}
	return false
}
