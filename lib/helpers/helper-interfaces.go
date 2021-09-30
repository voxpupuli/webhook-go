package helpers

import (
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

type Helpers interface {
	Normalize(allowUpper bool, str string) string
	GetPrefix(data parsers.Data, prefix string) string
	GetBranch(data parsers.Data, defaultBranch string) string
	GetEnvironment(branch, prefix string, allowUppercase bool) string
}
