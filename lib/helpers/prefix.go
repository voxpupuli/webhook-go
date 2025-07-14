package helpers

import (
	"fmt"
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
func (h *Helper) GetPrefixFromMapping(mapping map[string]string, repoName string) (string, error) {
	prefix, ok := mapping[repoName]
	if ok {
		return prefix, nil
	} else {
		return "", fmt.Errorf("Prefix not found in mapping for repo '%s'", repoName)
	}
}
