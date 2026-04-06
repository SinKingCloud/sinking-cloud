package enum

import (
	"service-admin/app/enum/menu_status"
	"service-admin/app/enum/menu_type"
	"service-admin/app/enum/permission"
	"service-admin/app/enum/role_name"
	"service-admin/app/enum/role_status"
	"service-admin/app/enum/user_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user_status.Status(), //用户状态
	},
	"role": func() interface{} {
		return map[string]interface{}{
			"status": role_status.Status(), //角色状态
			"names":  role_name.NameMap(),  //角色名称
		}
	},
	"menu": map[string]interface{}{
		"status": menu_status.Status(), //菜单状态
		"types":  menu_type.Types(),    //菜单类型
	},
	"permission": func() interface{} {
		return map[string]interface{}{
			"codes": permission.CodeMap(), //权限枚举
		}
	},
}
