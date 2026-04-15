package task

import (
	"service-cron/app/model"

	"github.com/robfig/cron/v3"
)

type Task struct {
	Info    *model.Task  //任务信息
	EntryID cron.EntryID //任务ID
}
