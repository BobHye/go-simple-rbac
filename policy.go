package go_simple_rbac

import (
	"encoding/json"
	"strings"
)

// A PermissionConstructor is a function that creates a new Permission from the specified action and target strings.
// PermissionConstructor 是一个函数，它根据指定的操作和目标字符串创建新的 Permission。
type PermissionConstructor func(action, target string) Permission

// DefaultPermissionConstructors returns a mapping of constructor names to PermissionConstructor functions
// 返回构造函数名称到 PermissionConstructor 函数的映射
// for each of the builtin PermissionConstructors:
//
//	"glob":   NewGlobPermission
//	"regex":  NewRegexPermission
//	"string": NewStringPermission
func DefaultPermissionConstructors() map[string]PermissionConstructor {
	return map[string]PermissionConstructor{
		"glob":   NewGlobPermission,
		"regex":  NewRegexPermission,
		"string": NewStringPermission,
	}
}

// A PermissionTemplate holds information about a permission in templated format.
// 以模板格式保存有关权限的信息。
type PermissionTemplate struct {
	Constructor string `json:"constructor,omitempty"`
	Action      string `json:"action"`
	Target      string `json:"target"`
}

// A PolicyTemplate holds information about a Role in a templated format.
// This format can be encoded to and from JSON.
// 以模板格式保存有关角色的信息。 此格式可以编码为 JSON 或从 JSON 编码。
type PolicyTemplate struct {
	RoleID              string               `json:"role_id"`
	PermissionTemplates []PermissionTemplate `json:"permissions"`
	constructors        map[string]PermissionConstructor
}

// NewPolicyTemplate generates a new PolicyTemplate with the specified roleID and default constructors.
// 使用指定的角色 ID 和默认构造函数生成一个新的 PolicyTemplate。
func NewPolicyTemplate(roleID string) *PolicyTemplate {
	return &PolicyTemplate{
		RoleID:              roleID,
		PermissionTemplates: []PermissionTemplate{},
		constructors:        DefaultPermissionConstructors(),
	}
}

// AddPermission adds a new PermissionTemplate to p.PermissionTemplates.
// 向 p.PermissionTemplates 添加一个新的 PermissionTemplate。
func (p *PolicyTemplate) AddPermission(constructor, action, target string) {
	p.PermissionTemplates = append(p.PermissionTemplates, PermissionTemplate{constructor, action, target})
}

// SetConstructor updates the mapping of a constructor name to a PermissionConstructor.
// If a mapping for the specified same name already exists, it will be overwritten.
// 将构造函数名称的映射更新为 PermissionConstructor。 如果指定的相同名称的映射已存在，则会被覆盖。
func (p *PolicyTemplate) SetConstructor(name string, constructor PermissionConstructor) {
	p.constructors[name] = constructor
}

// DeleteConstructor will remove the constructor mapping at the specified name if it exists.
// 如果存在，将删除指定名称的构造函数映射。
func (p *PolicyTemplate) DeleteConstructor(name string) {
	delete(p.constructors, name)
}

// Role converts the PolicyTemplate to a Role.
// Replacer can be used to replace variables within the Action and Target fields in the PermissionTemplates.
// Use GlobPermission as default if a PermissionTemplate.Constructor does not have a corresponding PermissionConstructor.
// Role 将 PolicyTemplate 转换为 Role。
// Replacer 可用于替换 PermissionTemplates 中的 Action 和 Target 字段内的变量。
// 如果 PermissionTemplate.Constructor 没有相应的 PermissionConstructor，则使用 GlobPermission 作为默认值。
func (p *PolicyTemplate) Role(replacer *strings.Replacer) *Role {
	role := &Role{
		RoleID:      p.RoleID,
		Permissions: make(Permissions, len(p.PermissionTemplates)),
	}

	for i, permissionTemplate := range p.PermissionTemplates {
		constructor, ok := p.constructors[permissionTemplate.Constructor]
		if !ok {
			constructor = NewGlobPermission
		}

		action := permissionTemplate.Action
		target := permissionTemplate.Target
		if replacer != nil {
			action = replacer.Replace(action)
			target = replacer.Replace(target)
		}
		role.Permissions[i] = constructor(action, target)
	}

	return role
}

// UnmarshalJSON allows a *PolicyTemplate to implement the json.Unmarshaler interface.
// We do this to set the default constructors on p after the unmarshalling.
// 允许 *PolicyTemplate 实现 json.Unmarshaler 接口。 我们这样做是为了在解组之后在 p 上设置默认构造函数。
func (p *PolicyTemplate) UnmarshalJSON(data []byte) error {
	type Alias PolicyTemplate
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.constructors = DefaultPermissionConstructors()
	return nil
}
