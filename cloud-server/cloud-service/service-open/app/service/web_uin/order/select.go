package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 查询
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.AppUinOrder, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppUinOrder{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["app_id"] != "" {
		query.Where("`app_id` = ?", where["app_id"])
	}
	if where["uin_id"] != "" {
		query.Where("`uin_id` = ?", where["uin_id"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
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
