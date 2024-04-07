package system_log_request

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取请求日志数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*LogRequest, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemLogRequest{})
	if where["request_id"] != "" {
		query.Where("`request_id` = ?", where["request_id"])
	}
	if where["request_ip"] != "" {
		query.Where("`request_ip` = ?", where["request_ip"])
	}
	if where["request_url"] != "" {
		query.Where("`request_url` like ?", where["request_url"]+"%")
	}
	if where["request_method"] != "" {
		query.Where("`request_method` = ?", where["request_method"])
	}
	if where["response_code"] != "" {
		query.Where("`response_code` = ?", where["response_code"])
	}
	if where["group_name"] != "" {
		query.Where("`group_name` = ?", where["group_name"])
	}
	if where["app_name"] != "" {
		query.Where("`app_name` = ?", where["app_name"])
	}
	if where["app_addr"] != "" {
		query.Where("`app_addr` = ?", where["app_addr"])
	}
	if where["response_time_max"] != "" {
		query.Where("`response_time` >= ?", where["response_time_max"])
	}
	if where["response_time_min"] != "" {
		query.Where("`response_time` <= ?", where["response_time_min"])
	}
	if where["request_header"] != "" {
		query.Where("`request_header` like ?", "%"+where["request_header"]+"%")
	}
	if where["request_body"] != "" {
		query.Where("`request_body` like ?", "%"+where["request_body"]+"%")
	}
	if where["response_header"] != "" {
		query.Where("`response_header` like ?", "%"+where["response_header"]+"%")
	}
	if where["response_body"] != "" {
		query.Where("`response_body` like ?", "%"+where["response_body"]+"%")
	}
	if where["request_time_start"] != "" {
		query.Where("`request_time` >= ?", where["request_time_start"])
	}
	if where["request_time_end"] != "" {
		query.Where("`request_time` <= ?", where["request_time_end"])
	}
	err = query.
		Count(&total).
		Limit(pageSize).
		Offset(offset).
		Order(orderByField + " " + orderByType).
		Find(&list).Error
	return
}
