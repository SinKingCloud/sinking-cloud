package task

import (
	"errors"
	"service-cron/app/model"
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// notify 通知其他节点重载数据
func (s *service) notify(action reload.Action, taskIds ...int64) {
	message := &reload.Message{
		Action:  action,
		Version: s.taskVersion(true),
	}
	uniqueMap := make(map[int64]struct{})
	for _, id := range taskIds {
		uniqueMap[id] = struct{}{}
	}
	ids := make([]int64, 0, len(uniqueMap))
	for id := range uniqueMap {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		message.TaskIds = ids
	}
	s.reloadQueue.SendBroadcastTask(message)
}

// log 记录任务日志
func (s *service) log(taskId int64, contents ...string) {
	content := &strings.Builder{}
	num := len(contents)
	if num > 0 {
		for i, v := range contents {
			content.WriteString(v)
			if i != len(contents)-1 {
				content.WriteString(" ")
			}
		}
	}
	_ = s.repositoryLog.Create(&model.Log{
		TaskId:     taskId,
		Content:    content.String(),
		CreateTime: str.DateTime(time.Now()),
	})
}

// ReloadQueueHandle 执行重载handle
func (s *service) ReloadQueueHandle(message *reload.Message) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	switch message.Action {
	case reload.Reload:
		for _, v := range message.TaskIds {
			info, err := s.FindById(v, false)
			if err != nil {
				continue
			}
			err = s.reload(info)
			if err != nil {
				s.log(info.Id, "载入任务失败", err.Error())
			}
		}
		if message.Version > s.localVersion.Load() {
			s.localVersion.Store(message.Version)
		}
	case reload.Remove:
		for _, v := range message.TaskIds {
			s.remove(v)
		}
		if message.Version > s.localVersion.Load() {
			s.localVersion.Store(message.Version)
		}
	case reload.Check:
		if message.Version > 0 || s.taskVersion(false) != s.localVersion.Load() {
			return s.ReloadAll()
		}
	}
	return nil
}

// TaskQueueHandle 执行任务handle
func (s *service) TaskQueueHandle(message *task.Message) error {
	return nil
}
