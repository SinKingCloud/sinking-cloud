package enum

import "service-admin/app/enum/user"

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user.Status(), //用户状态
	},
}
