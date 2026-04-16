package cron

import (
	"errors"
	"sync"
	"time"

	robfig "github.com/robfig/cron/v3"
)

// Task 定时任务
type Task struct {
	Name string
	Spec string
	Cmd  func()
}

// Info 定时任务详情
type Info struct {
	Name  string
	Spec  string
	Next  time.Time
	Prev  time.Time
	Valid bool
}

type entry struct {
	id   robfig.EntryID
	spec string
	cmd  func()
}

// Cron 定时任务实例
type Cron struct {
	ins     *robfig.Cron
	mu      sync.RWMutex
	parser  robfig.Parser
	entries sync.Map
}

// New 实例化定时任务
func New() *Cron {
	return &Cron{
		ins: robfig.New(
			robfig.WithSeconds(),
			robfig.WithChain(
				robfig.Recover(robfig.DiscardLogger),
			),
		),
		parser: robfig.NewParser(
			robfig.Second | robfig.Minute | robfig.Hour | robfig.Dom | robfig.Month | robfig.Dow | robfig.Descriptor,
		),
	}
}

// Validate 验证表达式
func (s *Cron) Validate(spec string) error {
	if spec == "" {
		return errors.New("任务表达式不能为空")
	}
	_, err := s.parser.Parse(spec)
	if err != nil {
		return errors.New("任务表达式验证失败")
	}
	return nil
}

// AddFunc 注册任务，同名时自动替换
func (s *Cron) AddFunc(name string, spec string, cmd func()) error {
	if name == "" {
		return errors.New("任务名称不能为空")
	}
	if cmd == nil {
		return errors.New("任务执行方法不能为空")
	}
	if err := s.Validate(spec); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.entries.Load(name); ok {
		if info, ok2 := value.(*entry); ok2 && info != nil && info.id > 0 {
			s.ins.Remove(info.id)
		}
	}
	id, err := s.ins.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	s.entries.Store(name, &entry{
		id:   id,
		spec: spec,
		cmd:  cmd,
	})
	return nil
}

// Add 注册任务对象
func (s *Cron) Add(task *Task) error {
	if task == nil {
		return errors.New("任务不能为空")
	}
	return s.AddFunc(task.Name, task.Spec, task.Cmd)
}

// Remove 移除任务
func (s *Cron) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.entries.Load(name); ok {
		if info, ok2 := value.(*entry); ok2 && info != nil && info.id > 0 {
			s.ins.Remove(info.id)
		}
		s.entries.Delete(name)
	}
}

// Detail 获取任务详情
func (s *Cron) Detail(name string) (*Info, error) {
	s.mu.RLock()
	value, ok := s.entries.Load(name)
	if !ok {
		s.mu.RUnlock()
		return nil, errors.New("任务不存在")
	}
	info, ok := value.(*entry)
	if !ok || info == nil || info.id <= 0 {
		s.mu.RUnlock()
		return nil, errors.New("任务不存在")
	}
	t := s.ins.Entry(info.id)
	s.mu.RUnlock()
	return &Info{
		Name:  name,
		Spec:  info.spec,
		Next:  t.Next,
		Prev:  t.Prev,
		Valid: t.Valid(),
	}, nil
}

// Run 主动执行任务
func (s *Cron) Run(name string, wait bool) error {
	s.mu.RLock()
	value, ok := s.entries.Load(name)
	if !ok {
		s.mu.RUnlock()
		return errors.New("任务不存在")
	}
	info, ok := value.(*entry)
	if !ok || info == nil || info.id <= 0 {
		s.mu.RUnlock()
		return errors.New("任务不存在")
	}
	t := s.ins.Entry(info.id)
	s.mu.RUnlock()
	if !t.Valid() || t.WrappedJob == nil {
		return errors.New("任务不存在")
	}
	if wait {
		t.WrappedJob.Run()
		return nil
	}
	go t.WrappedJob.Run()
	return nil
}

// Exists 任务是否存在
func (s *Cron) Exists(name string) bool {
	_, ok := s.entries.Load(name)
	return ok
}

// Start 启动
func (s *Cron) Start() {
	s.ins.Start()
}

// Stop 停止
func (s *Cron) Stop() {
	ctx := s.ins.Stop()
	<-ctx.Done()
}
