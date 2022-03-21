package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取请求日志数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.WebUserLog, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.WebUserLog{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["request_id"] != "" {
		query.Where("`request_id` = ?", where["request_id"])
	}
	if where["request_ip"] != "" {
		query.Where("`request_ip` = ?", where["request_ip"])
	}
	if where["title"] != "" {
		query.Where("`title` like ?", "%"+where["title"]+"%")
	}
	if where["content"] != "" {
		query.Where("`content` like ?", "%"+where["content"]+"%")
	}
	if where["create_time_start"] != "" {
		query.Where("`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("`create_time` <= ?", where["create_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
