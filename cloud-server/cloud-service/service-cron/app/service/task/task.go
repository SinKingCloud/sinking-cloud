package task

import (
	"service-cron/app/model"
	log2 "service-cron/app/queue/log"
	register2 "service-cron/app/queue/register"
	"service-cron/app/queue/reload"
	task2 "service-cron/app/queue/task"
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"
	"sync"
	"sync/atomic"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cron"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Service service接口
type Service interface {
	Close()
	FindById(id int64, refresh bool) (*model.Task, error)
	Select(where *task.SelectTask, page *page.Query) (list []*task.Task, total int64, err error)
	Create(data *model.Task) error
	DeleteByIds(ids []int64) error
	UpdateById(id int64, data *task.UpdateTask) error
	Run(id int64) error
	ReloadAll() error
	CheckReload() error
	Register(id int64, Type int, name string, spec string, script string, concurrent int, status int) error
}

// 任务服务
type service struct {
	repositoryLog  log.Interface
	repositoryTask task.Interface
	database       *database.Database
	cache          cache.Interface
	taskQueue      *queue.UnicastClient[*task2.Message]
	logQueue       *queue.UnicastClient[*log2.Message]
	reloadQueue    *queue.BroadcastClient[*reload.Message]
	registerQueue  *queue.UnicastClient[*register2.Message]
	cron           *cron.Cron
	entries        sync.Map
	localVersion   atomic.Int64
	dispatchNode   string
	isDispatchNode atomic.Bool
	started        atomic.Bool
	closed         atomic.Bool
	closeChan      chan struct{}
}

// NewService 实例化service
func NewService(
	repositoryLog log.Interface,
	repositoryTask task.Interface,
	cache cache.Interface,
	database *database.Database,
	taskQueue *queue.UnicastClient[*task2.Message],
	reloadQueue *queue.BroadcastClient[*reload.Message],
	logQueue *queue.UnicastClient[*log2.Message],
	registerQueue *queue.UnicastClient[*register2.Message],
) *service {
	ins := &service{
		repositoryLog:  repositoryLog,
		repositoryTask: repositoryTask,
		database:       database,
		cache:          cache,
		taskQueue:      taskQueue,
		reloadQueue:    reloadQueue,
		logQueue:       logQueue,
		registerQueue:  registerQueue,
		cron:           cron.New(),
		closeChan:      make(chan struct{}),
		dispatchNode:   "cron-" + str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5(true),
	}
	if taskQueue != nil {
		taskQueue.SetConsumer(ins.consumeTaskMessage)
	}
	if reloadQueue != nil {
		reloadQueue.SetConsumer(ins.consumeReloadMessage)
	}
	if logQueue != nil {
		logQueue.SetConsumer(ins.consumeLogMessage)
	}
	if registerQueue != nil {
		registerQueue.SetConsumer(ins.consumeRegisterMessage)
	}
	ins.start()
	return ins
}
