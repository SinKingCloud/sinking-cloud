package web_user_order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Delete 删除数据
func (s *Service) Delete(where map[string]string) (err error) {
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
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["trade_no"] != "" {
		query.Where("`trade_no` = ?", where["trade_no"])
	}
	if where["out_trade_no"] != "" {
		query.Where("`out_trade_no` = ?", where["out_trade_no"])
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
	err = query.Delete(&model.WebUserOrder{}).Error
	return
}
