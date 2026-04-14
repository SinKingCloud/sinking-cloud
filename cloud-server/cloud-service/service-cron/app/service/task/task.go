package task

import (
	"service-cron/app/repository/log"
	"service-cron/app/repository/task"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service 接口
type Service interface {
}

// Session服务
type service struct {
	repositoryLog  log.Interface
	repositoryTask task.Interface
	cache          cache.Interface
}

// NewService 实例化service
func NewService(repositoryLog log.Interface, repositoryTask task.Interface, cache cache.Interface) *service {
	return &service{
		repositoryLog:  repositoryLog,
		repositoryTask: repositoryTask,
		cache:          cache,
	}
}
