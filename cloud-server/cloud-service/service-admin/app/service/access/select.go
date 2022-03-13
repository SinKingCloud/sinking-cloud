package access

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取权限列表数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.SystemAccess, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemAccess{})
	if where["method"] != "" {
		query.Where("`method` = ?", where["method"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
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
	if where["url"] != "" {
		query.Where("`url` like ?", "%"+where["url"]+"%")
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
	}
	if where["describe"] != "" {
		query.Where("`describe` like ?", "%"+where["describe"]+"%")
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
