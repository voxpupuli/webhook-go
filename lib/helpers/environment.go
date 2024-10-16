package helpers

import (
	"fmt"
)

// GetEnvironment constructs and returns an environment name by combining the prefix and branch.
// If either is empty, it normalizes the branch name. Allows optional uppercase transformation.
func (h *Helper) GetEnvironment(branch, prefix string, allowUppercase bool) string {
	if prefix == "" || branch == "" {
		return h.Normalize(allowUppercase, branch)
	}
	return h.Normalize(allowUppercase, fmt.Sprintf("%s_%s", prefix, branch))
}
