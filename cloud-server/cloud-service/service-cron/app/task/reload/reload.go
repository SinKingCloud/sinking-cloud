package reload

import (
	"service-cron/app/queue/reload"
	"service-cron/app/service"
)

// Register 注册定时判断重载任务
func Register() (name string, spec string, cmd func()) {
	return "reload", "@every 1m", func() {
		_ = service.Task.ReloadQueueHandle(&reload.Message{
			Action: reload.Check,
		})
	}
}
