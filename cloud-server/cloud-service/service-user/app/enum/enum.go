package enum

import (
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/user_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user_status.Map(), //用户状态
	},
	"real_auth": map[string]interface{}{
		"status": real_auth_status.Map(), //实名状态
	},
}
