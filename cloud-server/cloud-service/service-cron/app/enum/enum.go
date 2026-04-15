package enum

import (
	"service-cron/app/enum/task_concurrent"
	"service-cron/app/enum/task_status"
	"service-cron/app/enum/task_type"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"task": map[string]interface{}{
		"status":      task_status.Map(),     //任务状态
		"types":       task_type.Map(),       //任务类型
		"concurrents": task_concurrent.Map(), //并发状态
	},
}
