package enum

import (
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/channel_type"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"channel": map[string]interface{}{
		"status": channel_status.Map(), //存储通道状态
		"type":   channel_type.Map(),   //存储通道类型
	},
	"file": map[string]interface{}{
		"status": channel_status.Map(), //文件状态
	},
	"scene": map[string]interface{}{
		"status": channel_status.Map(), //场景状态
	},
}
