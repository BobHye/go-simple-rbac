package go_simple_rbac

import "regexp"

// RegexMatch returns a Matcher that returns true
// if the target regular expression matches the specified pattern.
// 返回一个 Matcher，如果目标正则表达式与指定模式匹配，则返回 true
func RegexMatch(pattern string) Matcher {
	return func(target string) (bool, error) {
		return regexp.MatchString(pattern, target)
	}
}

// NewRegexPermission returns a Permission that uses RegexMatchers for the specified action and target patterns.
// 返回使用 RegexMatchers 执行指定操作和目标模式的权限。
func NewRegexPermission(actionPattern, targetPattern string) Permission {
	return NewPermission(RegexMatch(actionPattern), RegexMatch(targetPattern))
}
