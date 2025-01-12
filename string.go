package go_simple_rbac

// StringMatch returns a Matcher that returns true
// if the target string matches s.
// 返回一个 Matcher，如果目标字符串与 s 匹配，则返回 true
func StringMatch(s string) Matcher {
	return func(target string) (bool, error) {
		return target == s, nil
	}
}

// NewStringPermission returns a Permission that uses StringMatchers for the specified action and target.
// 返回使用 StringMatchers 执行指定操作和目标的权限
func NewStringPermission(action, target string) Permission {
	return NewPermission(StringMatch(action), StringMatch(target))
}
