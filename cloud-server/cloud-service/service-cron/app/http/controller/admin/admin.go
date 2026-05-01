package admin

import "service-cron/app/http/controller/admin/task"

type Controller struct {
	Task task.Controller
}
