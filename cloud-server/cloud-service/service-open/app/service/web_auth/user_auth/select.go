package user_auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 获取数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.UserAuth, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.UserAuth{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["auth_id"] != "" {
		query.Where("`auth_id` = ?", where["auth_id"])
	}
	if where["domain"] != "" {
		query.Where("`domain` = ?", where["domain"])
	}
	if where["code"] != "" {
		query.Where("`code` = ?", where["code"])
	}
	if where["ip"] != "" {
		query.Where("`ip` = ?", where["ip"])
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
