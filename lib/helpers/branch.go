package helpers

import (
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

// GetBranch returns the branch name from the parsed data. If the branch was deleted, it returns the defaultBranch.
func (h *Helper) GetBranch(data parsers.Data, defaultBranch string) string {
	if data.Deleted {
		return defaultBranch
	}
	return data.Branch
}
