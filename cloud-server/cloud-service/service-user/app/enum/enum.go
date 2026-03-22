package enum

import (
	"service-user/app/enum/real_auth"
	"service-user/app/enum/user"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user.Status(), //用户状态
	},
	"real_auth": map[string]interface{}{
		"status": real_auth.Status(), //实名状态
	},
}
