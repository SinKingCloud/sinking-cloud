package enum

import (
	"service-space/app/enum/mind_share_status"
	"service-space/app/enum/mind_type"
	"service-space/app/enum/mind_version_type"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"mind": map[string]interface{}{
		"type":         mind_type.Map(),         //预设分类
		"version_type": mind_version_type.Map(), //版本类型
		"share_status": mind_share_status.Map(), //分享状态
	},
}
