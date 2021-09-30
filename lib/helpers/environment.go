package helpers

import (
	"fmt"
)

func (h *Helper) GetEnvironment(branch, prefix string, allowUppercase bool) string {
	if prefix == "" || branch == "" {
		return h.Normalize(allowUppercase, branch)
	}
	return h.Normalize(allowUppercase, fmt.Sprintf("%s_%s", prefix, branch))
}
