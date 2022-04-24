package proxy

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// Select 查询
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.AppProxy, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppProxy{})
	if where["area_id"] != "" {
		query.Where("`area_id` = ?", where["area_id"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
