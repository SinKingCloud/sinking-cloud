package log

import (
	"service-cron/app/repository/log"
)

// Service 接口
type Service interface {
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
