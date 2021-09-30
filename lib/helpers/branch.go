package helpers

import (
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

func (h *Helper) GetBranch(data parsers.Data, defaultBranch string) string {
	if data.Deleted {
		return defaultBranch
	}
	return data.Branch
}
