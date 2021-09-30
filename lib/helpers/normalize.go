package helpers

import "strings"

func (h *Helper) Normalize(allowUpper bool, str string) string {
	if allowUpper {
		return str
	}
	return strings.ToLower(str)
}
