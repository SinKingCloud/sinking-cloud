package enum

import (
	"service-user/app/enum/message_status"
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/real_auth_type"
	"service-user/app/enum/user_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"user": map[string]interface{}{
		"status": user_status.Map(), //用户状态
	},
	"real_auth": map[string]interface{}{
		"status": real_auth_status.Map(), //实名状态
		"type":   real_auth_type.Map(),   //实名类型
	},
	"message": map[string]interface{}{
		"status": message_status.Map(), //消息状态
	},
}
