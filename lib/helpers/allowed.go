package helpers

import "path"

// BranchAllowed reports whether branch may be deployed according to the
// allowed_branches setting.
//
// Each entry is either an exact branch name or a glob pattern as understood by
// path.Match, e.g. "production", "release/*" or "feature_??". Note that "*" does
// not match "/", so use "release/*" to allow one level below "release/".
//
// An empty list means no restriction, so every branch is allowed; this keeps the
// default behaviour for existing configurations. When the list is not empty, a
// branch that matches no entry is not allowed. A malformed pattern never
// matches, so it fails closed. An empty branch is always allowed here because
// callers substitute the default branch before checking.
func (h *Helper) BranchAllowed(branch string, allowed []string) bool {
	if len(allowed) == 0 || branch == "" {
		return true
	}
	for _, pattern := range allowed {
		if pattern == "" {
			continue
		}
		if pattern == branch {
			return true
		}
		if ok, err := path.Match(pattern, branch); err == nil && ok {
			return true
		}
	}
	return false
}
