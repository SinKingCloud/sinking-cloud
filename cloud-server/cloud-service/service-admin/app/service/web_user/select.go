package web_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// SelectByIds 根据id查询用户
func (s *Service) SelectByIds(ids []int) (users []*model.WebUser, err error) {
	err = util.Database.Db.Model(&model.WebUser{}).Where("`id` in ? ", ids).Find(&users).Error
	return
}

// Select 查询用户
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*WebUser, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.WebUser{})
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["real_auth"] != "" {
		query.Where("`real_auth` = ?", where["real_auth"])
	}
	if where["account"] != "" {
		query.Where("`account` = ?", where["account"])
	}
	if where["email"] != "" {
		query.Where("`email` = ?", where["email"])
	}
	if where["phone"] != "" {
		query.Where("`phone` = ?", where["phone"])
	}
	if where["login_ip"] != "" {
		query.Where("`login_ip` = ?", where["login_ip"])
	}
	if where["nick_name"] != "" {
		query.Where("`nick_name` like ?", "%"+where["nick_name"]+"%")
	}
	if where["login_time_start"] != "" {
		query.Where("`login_time` >= ?", where["login_time_start"])
	}
	if where["login_time_end"] != "" {
		query.Where("`login_time` <= ?", where["login_time_end"])
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
