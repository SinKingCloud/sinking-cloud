package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// Count 统计
func (s *service) Count(where map[string]string) (total int64) {
	query := util.Database.Db.Model(&model.AppUinLog{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["app_id"] != "" {
		query.Where("`app_id` = ?", where["app_id"])
	}
	if where["uin_id"] != "" {
		query.Where("`uin_id` = ?", where["uin_id"])
	}
	if where["uin"] != "" {
		query.Where("`uin` = ?", where["uin"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	err := query.Count(&total).Error
	if err != nil {
		total = 0
	}
	return
}
