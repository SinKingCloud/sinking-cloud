package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 获取请求日志数据
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*SystemLog, total int64, err error) {
	offset := pageSize * (page - 1)
	table := &model.SystemUser{}
	query := util.Database.Db.Model(&model.SystemLog{}).
		Select("l.*,u.`account`,u.`nick_name`,u.`avatar`,u.`email`,u.`Phone`").
		Joins("AS l LEFT JOIN `" + table.TableName() + "` as u ON l.`user_id`=u.`id`")
	if where["request_id"] != "" {
		query.Where("l.`request_id` = ?", where["request_id"])
	}
	if where["request_ip"] != "" {
		query.Where("l.`request_ip` = ?", where["request_ip"])
	}
	if where["type"] != "" {
		query.Where("l.`type` = ?", where["type"])
	}
	if where["title"] != "" {
		query.Where("l.`title` like ?", "%"+where["title"]+"%")
	}
	if where["content"] != "" {
		query.Where("l.`content` like ?", "%"+where["content"]+"%")
	}
	if where["account"] != "" {
		query.Where("u.`account` like ?", "%"+where["account"]+"%")
	}
	if where["nick_name"] != "" {
		query.Where("u.`nick_name` like ?", "%"+where["nick_name"]+"%")
	}
	if where["email"] != "" {
		query.Where("u.`email` like ?", "%"+where["email"]+"%")
	}
	if where["phone"] != "" {
		query.Where("u.`phone` like ?", "%"+where["phone"]+"%")
	}
	if where["phone"] != "" {
		query.Where("u.`phone` like ?", "%"+where["phone"]+"%")
	}
	if where["create_time_start"] != "" {
		query.Where("l.`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("l.`create_time` <= ?", where["create_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Scan(&list).Error
	return
}
