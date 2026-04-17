package task

import (
	"encoding/json"
	"errors"
	"service-cron/app/constant"
	"service-cron/app/enum/task_concurrent"
	"service-cron/app/enum/task_status"
	"service-cron/app/enum/task_type"
	"service-cron/app/model"
	"service-cron/app/queue/reload"
	"service-cron/app/queue/task"
	repositoryTask "service-cron/app/repository/task"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Start 启动任务
func (s *service) Start() {
	_ = s.ReloadAll()
	s.cron.Start()
}

// Close 关闭任务
func (s *service) Close() {
	s.cron.Stop()
}

// taskVersion 获取或更新任务版本
func (s *service) taskVersion(set bool) int64 {
	cacheKey := constant.CacheNameWithTaskVersion
	if !set {
		value, exists := s.cache.Get(cacheKey)
		if !exists || value == "" {
			return 0
		}
		version, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0
		}
		return version
	}
	lockKey := cacheKey + "_lock"
	for i := 0; i < 100; i++ {
		version := int64(s.cache.Inc(cacheKey))
		if version > 0 {
			s.cache.SetExpire(cacheKey, constant.CacheTimeWithTaskVersion)
			return version
		}
		if s.cache.Lock(lockKey, 3000) {
			version = s.taskVersion(false) + 1
			s.cache.SetWithExpire(cacheKey, strconv.FormatInt(version, 10), constant.CacheTimeWithTaskVersion)
			s.cache.UnLock(lockKey)
			return version
		}
		time.Sleep(10 * time.Millisecond)
	}
	return 0
}

// FindById 根据ID查询任务
func (s *service) FindById(id int64, refresh bool) (data *model.Task, err error) {
	if !refresh {
		cacheKey := constant.CacheNameWithTaskInfo + strconv.FormatInt(id, 10)
		if value, exists := s.cache.Get(cacheKey); exists && value != "" {
			if err = json.Unmarshal([]byte(value), &data); err == nil && data != nil {
				return data, nil
			}
		}
	}
	data, err = s.repositoryTask.FindById(id)
	if err != nil {
		return nil, err
	}
	if value, err2 := json.Marshal(data); err2 == nil {
		s.cache.SetWithExpire(constant.CacheNameWithTaskInfo+strconv.FormatInt(id, 10), string(value), constant.CacheTimeWithTaskInfo)
	}
	return data, nil
}

// Select 查询任务列表
func (s *service) Select(where *repositoryTask.SelectTask, orderByField string, orderByType string, page int, pageSize int) (list []*repositoryTask.Task, total int64, err error) {
	return s.repositoryTask.Select(where, orderByField, orderByType, page, pageSize)
}

// Create 创建任务
func (s *service) Create(data *model.Task) (err error) {
	if data == nil {
		return errors.New("参数错误")
	}
	if _, ok := task_type.Map()[data.Type]; !ok {
		return errors.New("任务类型不合法")
	}
	if _, ok := task_status.Map()[data.Status]; !ok {
		return errors.New("任务状态不合法")
	}
	if _, ok := task_concurrent.Map()[data.Concurrent]; !ok {
		return errors.New("并发配置不合法")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("任务名称不能为空")
	}
	if strings.TrimSpace(data.Script) == "" {
		return errors.New("任务内容不能为空")
	}
	if strings.TrimSpace(data.Spec) == "" {
		return errors.New("任务表达式不能为空")
	}
	if err = s.cron.Validate(data.Spec); err != nil {
		return err
	}
	err = s.repositoryTask.Create(data)
	if err != nil || data.Id <= 0 {
		return errors.New("任务创建失败")
	}
	s.notify(reload.Reload, data.Id)
	return nil
}

// DeleteByIds 删除任务
func (s *service) DeleteByIds(ids []int64) error {
	if len(ids) <= 0 {
		return errors.New("删除数据不能为空")
	}
	err := s.database.Transaction(func(tx *gorm.DB) error {
		err := s.repositoryTask.DeleteByIds(ids, tx)
		if err != nil {
			return errors.Join(errors.New("删除任务数据失败"), err)
		}
		err = s.repositoryLog.DeleteByTaskIds(ids, tx)
		if err != nil {
			return errors.Join(errors.New("删除日志数据失败"), err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var keys []string
	for _, v := range ids {
		keys = append(keys, constant.CacheNameWithTaskInfo+strconv.FormatInt(v, 10))
	}
	s.cache.Delete(keys...)
	s.notify(reload.Remove, ids...)
	return nil
}

// UpdateById 通过ID更新
func (s *service) UpdateById(id int64, data *repositoryTask.UpdateTask) error {
	if id <= 0 || data == nil {
		return errors.New("更新数据不能为空")
	}
	if data.Type != nil {
		if value, ok := data.Type.(int); ok {
			if _, ok = task_type.Map()[value]; !ok {
				return errors.New("任务类型不合法")
			}
		} else {
			return errors.New("任务类型不合法")
		}
	}
	if data.Status != nil {
		if value, ok := data.Status.(int); ok {
			if _, ok = task_status.Map()[value]; !ok {
				return errors.New("任务状态不合法")
			}
		} else {
			return errors.New("任务状态不合法")
		}
	}
	if data.Concurrent != nil {
		if value, ok := data.Concurrent.(int); ok {
			if _, ok = task_concurrent.Map()[value]; !ok {
				return errors.New("并发配置不合法")
			}
		} else {
			return errors.New("并发配置不合法")
		}
	}
	if data.Name != nil {
		if value, ok := data.Name.(string); !ok || value == "" {
			return errors.New("任务名称不合法")
		}
	}
	if data.Script != nil {
		if value, ok := data.Script.(string); !ok || value == "" {
			return errors.New("任务内容不合法")
		}
	}
	if data.Spec != nil {
		value, ok := data.Spec.(string)
		if !ok || value == "" {
			return errors.New("任务表达式不合法")
		}
		if err := s.cron.Validate(value); err != nil {
			return err
		}
	}
	err := s.repositoryTask.UpdateById(id, data)
	if err != nil {
		return errors.Join(errors.New("修改任务失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithTaskInfo + strconv.FormatInt(id, 10))
	s.notify(reload.Reload, id)
	return nil
}

// Run 主动运行一次任务
func (s *service) Run(id int64) error {
	info, err := s.FindById(id, false)
	if err != nil {
		return err
	}
	if s.taskQueue != nil {
		s.taskQueue.SendUnicastTask(&task.Message{Id: info.Id})
	}
	return nil
}

// ReloadAll 重载所有任务
func (s *service) ReloadAll() error {
	targetVersion := s.taskVersion(false)
	list, err := s.repositoryTask.SelectAll()
	if err != nil {
		return errors.Join(errors.New("获取任务数据失败"), err)
	}
	taskIds := make(map[int64]struct{}, len(list))
	for _, v := range list {
		if v == nil || v.Id <= 0 {
			continue
		}
		taskIds[v.Id] = struct{}{}
		err = s.reload(v)
		if err != nil {
			s.log(v.Id, "载入任务失败", err.Error())
		}
	}
	var removeIds []int64
	s.entries.Range(func(key, value interface{}) bool {
		id, ok := key.(int64)
		if !ok {
			return true
		}
		if _, exists := taskIds[id]; !exists {
			removeIds = append(removeIds, id)
		}
		return true
	})
	for _, id := range removeIds {
		s.remove(id)
	}
	s.localVersion.Store(targetVersion)
	return nil
}

// remove 删除任务
func (s *service) remove(id int64) {
	s.cron.Remove(strconv.FormatInt(id, 10))
	s.entries.Delete(id)
}

// reload 重载任务
func (s *service) reload(data *model.Task) error {
	if data == nil || data.Id <= 0 || data.Spec == "" {
		return errors.New("任务不能为空")
	}
	if err := s.cron.Validate(data.Spec); err != nil {
		return err
	}
	cronName := strconv.FormatInt(data.Id, 10)
	specChange := false
	statusChange := false
	if value, ok := s.entries.Load(data.Id); ok && value != nil {
		if info, ok2 := value.(*model.Task); ok2 {
			if info != nil && info.Spec != data.Spec {
				specChange = true
			}
			if info != nil && info.Status != data.Status {
				statusChange = true
			}
		}
	} else {
		statusChange = data.Status == task_status.Running
	}
	s.entries.Store(data.Id, data)
	if data.Status == task_status.Stop {
		s.cron.Remove(cronName)
	} else if specChange || statusChange || !s.cron.Exists(cronName) {
		err := s.cron.AddFunc(cronName, data.Spec, func() {
			if s.taskQueue == nil {
				return
			}
			lockKey := constant.LockTaskDispatch + cronName
			if !s.cache.Lock(lockKey, constant.LockTimeTaskDispatch) {
				return
			}
			defer s.cache.UnLock(lockKey)
			s.taskQueue.SendUnicastTask(&task.Message{Id: data.Id})
		})
		if err != nil {
			return err
		}
	}
	return nil
}
