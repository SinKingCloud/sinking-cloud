package web_user_order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*WebUserOrder, total int64, err error) {
	offset := pageSize * (page - 1)
	table := &model.WebUser{}
	query := util.Database.Db.Model(&model.WebUserOrder{}).
		Select("o.*,u.`account`,u.`nick_name`,u.`avatar`,u.`email`,u.`phone`").
		Joins("AS o LEFT JOIN `" + table.TableName() + "` as u ON o.`user_id`=u.`id`")
	if where["user_id"] != "" {
		query.Where("o.`user_id` = ?", where["user_id"])
	}
	if where["pay_type"] != "" {
		query.Where("o.`pay_type` = ?", where["pay_type"])
	}
	if where["order_type"] != "" {
		query.Where("o.`order_type` = ?", where["order_type"])
	}
	if where["status"] != "" {
		query.Where("o.`status` = ?", where["status"])
	}
	if where["trade_no"] != "" {
		query.Where("o.`trade_no` = ?", where["trade_no"])
	}
	if where["out_trade_no"] != "" {
		query.Where("o.`out_trade_no` = ?", where["out_trade_no"])
	}
	if where["account"] != "" {
		query.Where("u.`account` = ?", where["account"])
	}
	if where["email"] != "" {
		query.Where("u.`email` = ?", where["email"])
	}
	if where["phone"] != "" {
		query.Where("u.`phone` = ?", where["phone"])
	}
	if where["nick_name"] != "" {
		query.Where("u.`nick_name` like ?", "%"+where["nick_name"]+"%")
	}
	if where["name"] != "" {
		query.Where("u.`name` like ?", "%"+where["name"]+"%")
	}
	if where["create_time_start"] != "" {
		query.Where("o.`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("o.`create_time` <= ?", where["create_time_end"])
	}
	if where["update_time_start"] != "" {
		query.Where("o.`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("o.`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Scan(&list).Error
	return
}
