package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// Count 统计
func (s *Service) Count(where map[string]string) (total int64) {
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
	if where["trade_no"] != "" {
		query.Where("`trade_no` = ?", where["trade_no"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
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
	err := query.Count(&total).Error
	if err != nil {
		total = 0
	}
	return
}
