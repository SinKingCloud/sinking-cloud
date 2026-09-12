package log

import (
	"service-cron/app/model"
	"service-cron/app/repository/log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectByTaskId 查询数据
func (s *service) SelectByTaskId(id int64, query *page.Query) (*page.Result[*model.Log], error) {
	where := &log.SelectLog{
		TaskId: &id,
	}
	return s.repositoryLog.Select(where, query)
}
