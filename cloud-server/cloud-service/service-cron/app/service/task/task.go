package task

import (
	"service-cron/app/model"
	"service-cron/app/queue/reload"
	task2 "service-cron/app/queue/task"
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cron"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// Service service接口
type Service interface {
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
}

// NewService 实例化service
func NewService(repositoryLog log.Interface, repositoryTask task.Interface, cache cache.Interface, database *database.Database, taskQueue *queue.UnicastClient[*task2.Message], reloadQueue *queue.BroadcastClient[*reload.Message]) *service {
	ins := &service{
		repositoryLog:  repositoryLog,
		repositoryTask: repositoryTask,
		database:       database,
		cache:          cache,
		taskQueue:      taskQueue,
		reloadQueue:    reloadQueue,
		cron:           cron.New(),
	}
	ins.Start()
	return ins
}
