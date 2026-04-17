package task

import (
	"service-cron/app/constant"
	"service-cron/app/model"
	"service-cron/app/queue/reload"
	task2 "service-cron/app/queue/task"
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cron"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Service service接口
type Service interface {
	SetQueue(taskQueue *queue.UnicastClient[*task2.Message], reloadQueue *queue.BroadcastClient[*reload.Message])
	Start()
	Close()
	FindById(id int64, refresh bool) (*model.Task, error)
	Select(where *task.SelectTask, orderByField string, orderByType string, page int, pageSize int) (list []*task.Task, total int64, err error)
	Create(data *model.Task) error
	DeleteByIds(ids []int64) error
	UpdateById(id int64, data *task.UpdateTask) error
	Run(id int64) error
	ReloadAll() error
	TaskQueueHandle(message *task2.Message) error
	ReloadQueueHandle(message *reload.Message) error
}

// Session服务
type service struct {
	repositoryLog  log.Interface
	repositoryTask task.Interface
	database       *database.Database
	cache          cache.Interface
	taskQueue      *queue.UnicastClient[*task2.Message]
	reloadQueue    *queue.BroadcastClient[*reload.Message]
	cron           *cron.Cron
	entries        sync.Map
	localVersion   atomic.Int64
	dispatchNode   string
	isDispatchNode atomic.Bool
	started        atomic.Bool
	closed         atomic.Bool
	closeChan      chan struct{}
}

// SetQueue 设置队列实例
func (s *service) SetQueue(taskQueue *queue.UnicastClient[*task2.Message], reloadQueue *queue.BroadcastClient[*reload.Message]) {
	s.taskQueue = taskQueue
	s.reloadQueue = reloadQueue
}

// NewService 实例化service
func NewService(repositoryLog log.Interface, repositoryTask task.Interface, cache cache.Interface, database *database.Database) *service {
	ins := &service{
		repositoryLog:  repositoryLog,
		repositoryTask: repositoryTask,
		database:       database,
		cache:          cache,
		cron:           cron.New(),
		closeChan:      make(chan struct{}),
		dispatchNode:   "cron-" + str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5(true),
	}
	return ins
}

// keepDispatchNode 保持投递节点身份
func (s *service) keepDispatchNode() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.closeChan:
			if s.isDispatchNode.Load() {
				if value, exists := s.cache.Get(constant.CacheNameWithTaskDispatchNode); exists && value == s.dispatchNode {
					s.cache.Delete(constant.CacheNameWithTaskDispatchNode)
				}
				s.isDispatchNode.Store(false)
			}
			return
		case <-ticker.C:
			if s.isDispatchNode.Load() {
				if !s.cache.Lock(constant.LockTaskDispatchRenew, constant.LockTimeTaskDispatchRenew) {
					continue
				}
				value, exists := s.cache.Get(constant.CacheNameWithTaskDispatchNode)
				if !exists || value == "" || value != s.dispatchNode {
					s.isDispatchNode.Store(false)
					s.cache.UnLock(constant.LockTaskDispatchRenew)
					continue
				}
				if !s.cache.SetWithExpire(constant.CacheNameWithTaskDispatchNode, s.dispatchNode, constant.CacheTimeWithTaskDispatchNode) {
					s.isDispatchNode.Store(false)
				}
				s.cache.UnLock(constant.LockTaskDispatchRenew)
				continue
			}
			if value, exists := s.cache.Get(constant.CacheNameWithTaskDispatchNode); exists && value != "" {
				if value == s.dispatchNode {
					if !s.isDispatchNode.Load() {
						s.isDispatchNode.Store(true)
					}
				} else if s.isDispatchNode.Load() {
					s.isDispatchNode.Store(false)
				}
				continue
			}
			if !s.cache.Lock(constant.LockTaskDispatch, constant.LockTimeTaskDispatch) {
				continue
			}
			if value, exists := s.cache.Get(constant.CacheNameWithTaskDispatchNode); exists && value != "" {
				if value == s.dispatchNode {
					if !s.isDispatchNode.Load() {
						s.isDispatchNode.Store(true)
					}
				} else if s.isDispatchNode.Load() {
					s.isDispatchNode.Store(false)
				}
				s.cache.UnLock(constant.LockTaskDispatch)
				continue
			}
			if s.cache.SetWithExpire(constant.CacheNameWithTaskDispatchNode, s.dispatchNode, constant.CacheTimeWithTaskDispatchNode) {
				s.isDispatchNode.Store(true)
			}
			s.cache.UnLock(constant.LockTaskDispatch)
		}
	}
}
