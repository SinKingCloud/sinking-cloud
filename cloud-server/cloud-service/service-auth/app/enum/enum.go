package enum

import "service-auth/app/enum/session"

// Data 枚举信息
var Data = map[string]interface{}{
	"session": map[string]interface{}{
		"user_type": session.UserTypes(), //用户类型数据
	},
}
