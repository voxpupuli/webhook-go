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

// Lookup the repo name in the mapping and return the corresponding prefix.
// If the repo name is not found in the mapping, it returns an empty string.
func (h *Helper) GetPrefixFromMapping(mapping map[string]string, repoName string) string {
	return mapping[repoName]
}
