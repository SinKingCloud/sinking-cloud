package rbac

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// MenuTree 菜单树
type MenuTree struct {
	Id       int64       `json:"id"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Icon     string      `json:"icon"`
	Sort     int         `json:"sort"`
	Children []*MenuTree `json:"children,omitempty"`
}

// Role 角色详情
type Role struct {
	Id         int64        `json:"id"`
	Name       string       `json:"name"`
	Status     int          `json:"status"`
	MenuIds    []int64      `json:"menu_ids"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Control 页面控件
type Control struct {
	Id              int64        `json:"id"`
	Code            string       `json:"code"`
	Name            string       `json:"name"`
	Icon            string       `json:"icon"`
	PermissionCodes []string     `json:"permission_codes"`
	Status          int          `json:"status"`
	CreateTime      str.DateTime `json:"create_time"`
	UpdateTime      str.DateTime `json:"update_time"`
}

// UserMenuTree 用户菜单树
type UserMenuTree struct {
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Type     int             `json:"type"`
	Icon     string          `json:"icon"`
	Sort     int             `json:"sort"`
	Children []*UserMenuTree `json:"children,omitempty"`
}

// UserControl 用户页面控件
type UserControl struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// MenuOrControlInfo 菜单或控件详情
type MenuOrControlInfo struct {
	Id              int64        `json:"id"`
	Code            string       `json:"code"`
	Name            string       `json:"name"`
	Icon            string       `json:"icon"`
	Type            int          `json:"type"`
	Sort            int          `json:"sort"`
	Status          int          `json:"status"`
	PermissionCodes []string     `json:"permission_codes,omitempty"`
	CreateTime      str.DateTime `json:"create_time"`
	UpdateTime      str.DateTime `json:"update_time"`
}
