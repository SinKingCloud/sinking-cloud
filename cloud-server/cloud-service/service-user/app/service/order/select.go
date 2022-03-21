package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取请求日志数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.WebUserOrder, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.WebUserOrder{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["pay_type"] != "" {
		query.Where("`pay_type` = ?", where["pay_type"])
	}
	if where["order_type"] != "" {
		query.Where("`order_type` = ?", where["order_type"])
	}
	if where["trade_no"] != "" {
		query.Where("`trade_no` = ?", where["trade_no"])
	}
	if where["out_trade_no"] != "" {
		query.Where("`out_trade_no` = ?", where["out_trade_no"])
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
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
