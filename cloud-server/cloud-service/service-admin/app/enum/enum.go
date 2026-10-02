package enum

import (
	"service-admin/app/enum/menu_status"
	"service-admin/app/enum/menu_type"
	"service-admin/app/enum/permission_code"
	"service-admin/app/enum/role_name"
	"service-admin/app/enum/role_status"
	"service-admin/app/enum/user_log_type"
	"service-admin/app/enum/user_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user_status.Map(), //用户状态
	},
	"user_log": map[string]interface{}{
		"type": user_log_type.Map(), //操作日志类型
	},
	"role": func() interface{} {
		return map[string]interface{}{
			"status": role_status.Map(), //角色状态
			"names":  role_name.Map(),   //角色名称
		}
	},
	"menu": map[string]interface{}{
		"status": menu_status.Map(), //菜单状态
		"types":  menu_type.Map(),   //菜单类型
	},
	"permission": func() interface{} {
		return map[string]interface{}{
			"codes": permission_code.Map(), //权限枚举
		}
	},
}
