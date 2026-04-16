package task

import (
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"
)

// notify 通知其他节点重载数据
func (s *service) notify(action reload.Action, taskIds ...int64) {
	message := &reload.Message{
		Action: action,
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

// ReloadQueueHandle 执行重载handle
func (s *service) ReloadQueueHandle(message *reload.Message) error {
	switch message.Action {
	case reload.Reload:
		for _, v := range message.TaskIds {
			info, err := s.FindById(v, false)
			if err == nil && info != nil && info.Id > 0 {
				_ = s.reload(info)
			}
		}
		break
	case reload.Remove:
		for _, v := range message.TaskIds {
			s.remove(v)
		}
		break
	case reload.Check:

		break
	}
	return nil
}

// TaskQueueHandle 执行任务handle
func (s *service) TaskQueueHandle(message *task.Message) error {
	return nil
}
