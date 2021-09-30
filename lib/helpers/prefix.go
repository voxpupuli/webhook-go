package helpers

import (
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

func (h *Helper) GetPrefix(data parsers.Data, prefix string) string {
	switch prefix {
	case "repo":
		return data.RepoName
	case "user":
		return data.RepoUser
	case "":
		return ""
	default:
		return prefix
	}
}
