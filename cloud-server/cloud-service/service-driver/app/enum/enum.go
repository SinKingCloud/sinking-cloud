package enum

import (
	"service-driver/app/enum/channel_name"
	"service-driver/app/enum/channel_status"
	"service-driver/app/enum/channel_type"
	"service-driver/app/enum/file_status"
	"service-driver/app/enum/scene_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"channel": func() interface{} {
		return map[string]interface{}{
			"status": channel_status.Map(), //存储通道状态
			"type":   channel_type.Map(),   //存储通道类型
			"name":   channel_name.Map(),   //存储通道名称
		}
	},
	"file": map[string]interface{}{
		"status": file_status.Map(), //文件状态
	},
	"scene": func() interface{} {
		return map[string]interface{}{
			"status": scene_status.Map(), //场景状态
			"name":   scene_status.Map(), //场景名称
		}
	},
}
