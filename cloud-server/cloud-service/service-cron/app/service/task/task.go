package task

import (
	"service-cron/app/model"
	"service-cron/app/queue/reload"
	task2 "service-cron/app/queue/task"
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/robfig/cron/v3"
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
	ReloadAll() error
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
	parser         cron.Parser
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
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.Recover(cron.DiscardLogger),
			),
		),
		parser: cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		),
	}
	ins.Start()
	return ins
}
