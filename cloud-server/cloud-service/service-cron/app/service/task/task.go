package task

import (
	"service-cron/app/model"
	log2 "service-cron/app/queue/log"
	"service-cron/app/queue/reload"
	task2 "service-cron/app/queue/task"
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"
	"sync"
	"sync/atomic"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cron"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Service service接口
type Service interface {
	SetQueue(taskQueue *queue.UnicastClient[*task2.Message], reloadQueue *queue.BroadcastClient[*reload.Message], logQueue *queue.UnicastClient[*log2.Message])
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
	LogQueueHandle(message *log2.Message) error
	Register(id int64, Type int, name string, spec string, script string, concurrent int, status int) error
}

// Session服务
type service struct {
	repositoryLog  log.Interface
	repositoryTask task.Interface
	database       *database.Database
	cache          cache.Interface
	taskQueue      *queue.UnicastClient[*task2.Message]
	logQueue       *queue.UnicastClient[*log2.Message]
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
func (s *service) SetQueue(taskQueue *queue.UnicastClient[*task2.Message], reloadQueue *queue.BroadcastClient[*reload.Message], logQueue *queue.UnicastClient[*log2.Message]) {
	s.taskQueue = taskQueue
	s.reloadQueue = reloadQueue
	s.logQueue = logQueue
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
