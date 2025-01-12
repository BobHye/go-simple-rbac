package go_simple_rbac

import "strings"

// GLOB The character which is treated like a glob
const GLOB = "*"

// Glob will test a string pattern, potentially containing globs, against a
// subject string. The result is a simple true/false, determining whether or
// not the glob pattern matched the subject text.
func Glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// if the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)
	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])
		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}
		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}
	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}

// GlobMatch returns a Matcher that returns true if the target glob matches the specified pattern.
// GlobMatch 返回一个 Matcher，如果目标 glob 与指定模式匹配，则返回 true
func GlobMatch(pattern string) Matcher {
	return func(target string) (bool, error) {
		return Glob(pattern, target), nil
	}
}

// NewGlobPermission returns a Permission that uses GlobMatchers for the specified action and target patterns.
// 返回使用 GlobMatchers 指定操作和目标模式的权限
func NewGlobPermission(actionPattern, targetPattern string) Permission {
	return NewPermission(GlobMatch(actionPattern), GlobMatch(targetPattern))
}
