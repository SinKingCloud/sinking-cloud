package log

import (
	"service-cron/app/model"
	"service-cron/app/repository/log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Service 接口
type Service interface {
	SelectByTaskId(id int64, page *page.Query) (list []*model.Log, total int64, err error)
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
