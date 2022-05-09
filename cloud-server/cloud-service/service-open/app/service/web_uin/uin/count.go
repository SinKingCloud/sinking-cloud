package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Count 统计
func (s *service) Count(where map[string]string) (total int64) {
	query := util.Database.Db.Model(&model.AppUin{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["app_id"] != "" {
		query.Where("`app_id` = ?", where["app_id"])
	}
	if where["server_id"] != "" {
		query.Where("`server_id` = ?", where["server_id"])
	}
	if where["server_key"] != "" {
		query.Where("`server_key` = ?", where["server_key"])
	}
	if where["uin_status"] != "" {
		query.Where("`uin_status` = ?", where["uin_status"])
	}
	if where["uin"] != "" {
		query.Where("`uin` = ?", where["uin"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["run_time_start"] != "" {
		query.Where("`run_time` >= ?", where["run_time_start"])
	}
	if where["run_time_end"] != "" {
		query.Where("`run_time` <= ?", where["run_time_end"])
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
	err := query.Count(&total).Error
	if err != nil {
		total = 0
	}
	return
}
