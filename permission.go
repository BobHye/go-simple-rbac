package go_simple_rbac

// A Permission is a function that returns true if the action is allowed on the target
// 权限是一个函数，如果目标上允许执行该操作，则返回 true
type Permission func(action string, target string) (bool, error)

// The Permissions type is an adapter to allow helper functions to execute on a slice of Permissions
// Permissions 类型是一个适配器，允许辅助函数在 Permissions 片段上执行
type Permissions []Permission

// Can returns true if at least one of the permissions in p allows the action on the target
// 如果 p 中至少有一个权限允许对目标执行操作，则返回 true
func (p Permissions) Can(action string, target string) (bool, error) {
	for _, permission := range p {
		can, err := permission(action, target)
		if err != nil {
			return false, err
		}

		if can {
			return true, nil
		}
	}

	return false, nil
}

// NewPermission returns a Permission that will return true | 返回将返回 true 的 Permission
// if the actionMatcher returns true for the given action, and
// if the targetMatcher returns true the given target.
// 如果 actionMatcher 对给定的动作返回 true，并且如果 targetMatcher 对给定的目标返回 true。
func NewPermission(actionMatcher, targetMatcher Matcher) Permission {
	return func(action string, target string) (bool, error) {
		actionMatch, err := actionMatcher(action)
		if err != nil {
			return false, err
		}

		if !actionMatch {
			return false, nil
		}

		return targetMatcher(target)
	}
}

// AllowAll is a Permission that always returns true
// AllowAll 是一个始终返回 true 的权限
func AllowAll(action, target string) (bool, error) {
	return true, nil
}
