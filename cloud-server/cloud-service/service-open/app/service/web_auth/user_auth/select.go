package user_auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*UserAuth, total int64, err error) {
	offset := pageSize * (page - 1)
	table := &model.AuthApp{}
	query := util.Database.Db.Model(&model.UserAuth{}).Select("u.*,a.`name`,a.`describe`,a.`version_id`,a.`release_down_num`,a.`update_down_num`,a.`price`").
		Joins("AS u LEFT JOIN `" + table.TableName() + "` as a ON u.`auth_id`=a.`id`")
	if where["user_id"] != "" {
		query.Where("u.`user_id` = ?", where["user_id"])
	}
	if where["auth_id"] != "" {
		query.Where("u.`auth_id` = ?", where["auth_id"])
	}
	if where["domain"] != "" {
		query.Where("u.`domain` = ?", where["domain"])
	}
	if where["code"] != "" {
		query.Where("u.`code` = ?", where["code"])
	}
	if where["ip"] != "" {
		query.Where("u.`ip` = ?", where["ip"])
	}
	if where["create_time_start"] != "" {
		query.Where("u.`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("u.`create_time` <= ?", where["create_time_end"])
	}
	if where["update_time_start"] != "" {
		query.Where("u.`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("u.`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Scan(&list).Error
	return
}
