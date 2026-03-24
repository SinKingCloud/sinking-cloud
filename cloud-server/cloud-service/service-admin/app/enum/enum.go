package enum

import (
	"service-admin/app/enum/role"
	"service-admin/app/enum/user"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user.Status(), //用户状态
	},
	"role": map[string]interface{}{
		"status": role.Status(), //角色状态
	},
}
