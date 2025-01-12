package go_simple_rbac

// A Matcher is a function that returns a bool representing
// whether or not the target matches some pre-defined pattern.
type Matcher func(target string) (bool, error)

// MatchAny will convert a slice of Matchers into a single Matcher
// that returns true if and only if at least one of the specified matchers returns true.
// Matcher 是一个返回布尔值的函数，表示目标是否与某些预定义模式匹配
func MatchAny(matchers ...Matcher) Matcher {
	return func(target string) (bool, error) {
		for _, matcher := range matchers {
			match, err := matcher(target)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
		}
		return false, nil
	}
}

// MatchAll will convert a slice of Matchers into a single Matcher
// that returns true if and only if all of the specified matchers returns true.
// MatchAll 会将 Matcher 的一个片段转换为单个 Matcher 当且仅当所有指定的匹配器都返回 true 时，才会返回 true。
func MatchAll(matchers ...Matcher) Matcher {
	return func(target string) (bool, error) {
		for _, matcher := range matchers {
			match, err := matcher(target)
			if err != nil {
				return false, err
			}

			if !match {
				return false, nil
			}
		}

		return true, nil
	}
}

// Anything is a Matcher that always returns true
// 任何东西都是始终返回 true 的 Matcher
func Anything(target string) (bool, error) {
	return true, nil
}
