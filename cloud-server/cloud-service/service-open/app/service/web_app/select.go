package web_app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 获取请求日志数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.UserApp, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.UserApp{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["key"] != "" {
		query.Where("`key` = ?", where["key"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
	}
	if where["describe"] != "" {
		query.Where("`describe` like ?", "%"+where["describe"]+"%")
	}
	if where["create_time_start"] != "" {
		query.Where("`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("`create_time` <= ?", where["create_time_end"])
	}
	if where["update_time_start"] != "" {
		query.Where("`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
