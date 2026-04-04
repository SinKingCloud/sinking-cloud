package rbac

// MenuTree 菜单树
type MenuTree struct {
	Id       int64       `json:"-"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Icon     string      `json:"icon"`
	Sort     int         `json:"sort"`
	Children []*MenuTree `json:"children,omitempty"`
}

// Control 页面控件
type Control struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
