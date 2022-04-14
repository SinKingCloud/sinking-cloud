package web_notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*WebNotice, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.WebNotice{})
	if where["place"] != "" {
		query.Where("`place` = ?", where["place"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["sort"] != "" {
		query.Where("`sort` = ?", where["sort"])
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
	if where["update_time_start"] != "" {
		query.Where("`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
