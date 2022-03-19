package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 查询角色
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.SystemConfig, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemConfig{})
	if where["pid"] != "" {
		query.Where("`pid` = ?", where["pid"])
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
	}
	if where["describe"] != "" {
		query.Where("`describe` like ?", "%"+where["describe"]+"%")
	}
	if where["key"] != "" {
		query.Where("`key` like ?", "%"+where["key"]+"%")
	}
	if where["value"] != "" {
		query.Where("`value` like ?", "%"+where["value"]+"%")
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
