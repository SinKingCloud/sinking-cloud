package rbac

// UserMenu 用户菜单
type UserMenu struct {
	Id       int64       `json:"-"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Icon     string      `json:"icon"`
	Sort     int         `json:"sort"`
	Children []*UserMenu `json:"children,omitempty"`
}

// UserControl 用户页面控件
type UserControl struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
