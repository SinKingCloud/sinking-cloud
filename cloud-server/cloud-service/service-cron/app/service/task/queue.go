package task

import (
	"errors"
	"service-cron/app/constant"
	"service-cron/app/enum/task_concurrent"
	"service-cron/app/enum/task_type"
	"service-cron/app/model"
	"service-cron/app/queue/log"
	"service-cron/app/queue/register"
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"
	"strconv"
	"strings"
	"time"

	commonConstant "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cmd"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
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
	s.logQueue.SendUnicastTask(&log.Message{
		Log: &model.Log{
			Id:         str.GetSnowWorkIns().GetId(),
			TaskId:     taskId,
			Content:    content.String(),
			CreateTime: str.DateTime(time.Now()),
		},
	})
}

// LogQueueHandle 记录日志
func (s *service) LogQueueHandle(message *log.Message) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	return s.repositoryLog.Create(message.Log)
}

// RegisterQueueHandle 注册任务handle
func (s *service) RegisterQueueHandle(message *register.Message) error {
	if message == nil {
		return errors.New("注册任务消息不能为空")
	}
	err := s.register(message.Id, message.Type, message.Name, message.Spec, message.Script, message.Concurrent, message.Status)
	if err != nil {
		s.log(message.Id, "注册任务失败", err.Error())
	}
	return err
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
	if message == nil || message.Id <= 0 {
		return errors.New("任务消息不能为空")
	}
	var (
		info *model.Task
		err  error
	)
	if value, ok := s.entries.Load(message.Id); ok && value != nil {
		if data, ok2 := value.(*model.Task); ok2 && data != nil {
			info = data
		}
	}
	if info == nil {
		info, err = s.FindById(message.Id, false)
		if err != nil {
			s.log(message.Id, "获取任务失败", err.Error())
			return err
		}
	}
	if info == nil || info.Id <= 0 {
		return errors.New("任务不存在")
	}
	if strings.TrimSpace(info.Script) == "" {
		err = errors.New("任务内容不能为空")
		s.log(info.Id, err.Error())
		return err
	}
	if info.Concurrent == task_concurrent.Disable {
		lockKey := constant.LockTask + strconv.FormatInt(info.Id, 10)
		lockDeadline := time.Now().Add(time.Duration(constant.LockTimeTask) * time.Millisecond)
		locked := false
		for time.Now().Before(lockDeadline) {
			if s.cache.Lock(lockKey, constant.LockTimeTask) {
				locked = true
				break
			}
			time.Sleep(time.Second)
		}
		if !locked {
			err = errors.New("任务执行超时")
			s.log(info.Id, err.Error())
			return err
		}
		defer s.cache.UnLock(lockKey)
	}
	s.log(info.Id, "开始执行任务", info.Name)
	switch info.Type {
	case task_type.Script:
		exec := cmd.NewScriptExec(commonConstant.TempPath, 600, func(line string) {
			if strings.TrimSpace(line) != "" {
				s.log(info.Id, line)
			}
		})
		_, stderr, err2 := exec.Execute(info.Script)
		if strings.TrimSpace(stderr) != "" {
			s.log(info.Id, stderr)
		}
		if err2 != nil {
			s.log(info.Id, "执行脚本失败", err2.Error())
			return err2
		}
	case task_type.Rpc:
		serviceName, action, ok := strings.Cut(info.Script, ".")
		serviceName = strings.TrimSpace(serviceName)
		action = strings.TrimSpace(action)
		if !ok || serviceName == "" || action == "" {
			err = errors.New("rpc任务格式字段不合法")
			s.log(info.Id, err.Error())
			return err
		}
		err = rpc.InvokeFunc(serviceName, action, nil, nil)
		if err != nil {
			s.log(info.Id, "执行RPC失败", err.Error())
			return err
		}
	default:
		err = errors.New("任务类型字段不合法")
		s.log(info.Id, err.Error())
		return err
	}
	s.log(info.Id, "任务执行成功", info.Name)
	return nil
}
