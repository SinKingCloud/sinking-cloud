package log

import (
	"service-cron/app/model"
	"service-cron/app/repository/log"
	"strconv"
)

// SelectByTaskId 查询数据
func (s *service) SelectByTaskId(id int64, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error) {
	where := &log.SelectLog{
		TaskId: strconv.FormatInt(id, 10),
	}
	return s.repositoryLog.Select(where, orderByField, orderByType, page, pageSize)
}
