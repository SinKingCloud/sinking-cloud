package log

import (
	"service-cron/app/model"
	"service-cron/app/repository/log"
)

// Service 接口
type Service interface {
	SelectByTaskId(id int64, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error)
}

// Session服务
type service struct {
	repositoryLog log.Interface
}

// NewService 实例化service
func NewService(repositoryLog log.Interface) *service {
	return &service{
		repositoryLog: repositoryLog,
	}
}
