package helpers

import (
	"regexp"
	"strings"
)

// BranchMatches reports whether branch matches any entry of patterns. It is the
// matching used by both allowed_branches and blocked_branches, so the two
// settings accept the same syntax.
//
// An entry is an exact branch name, optionally containing wildcards: "*" matches
// any sequence of characters (including "/") and "?" matches exactly one. No
// other character is special; git does not allow "*", "?" or "[" in branch
// names, so an existing exact entry can never be read as a pattern. Empty
// entries are ignored.
func (h *Helper) BranchMatches(branch string, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		if pattern == branch || globMatch(pattern, branch) {
			return true
		}
	}
	return false
}

// BranchAllowed reports whether branch may be deployed according to the
// allowed_branches setting. An empty list means no restriction, so every branch
// is allowed; this keeps the default behaviour for existing configurations.
// When the list is not empty, a branch that matches no entry is not allowed. An
// empty branch is always allowed here because callers substitute the default
// branch before checking.
func (h *Helper) BranchAllowed(branch string, allowed []string) bool {
	if len(allowed) == 0 || branch == "" {
		return true
	}
	return h.BranchMatches(branch, allowed)
}

func globMatch(pattern, s string) bool {
	if !strings.ContainsAny(pattern, "*?") {
		return false
	}
	var b strings.Builder
	b.WriteString("^")
	for _, r := range pattern {
		switch r {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteString(".")
		default:
			b.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile("(?s)" + b.String())
	return err == nil && re.MatchString(s)
}
