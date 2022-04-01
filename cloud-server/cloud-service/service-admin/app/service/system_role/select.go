package system_role

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// SelectByIds 根据id查询角色
func (s *service) SelectByIds(ids []int) (roles []*model.SystemRole, err error) {
	err = util.Database.Db.Model(&model.SystemRole{}).Where("`id` in ? ", ids).Find(&roles).Error
	return
}

// SelectAll 根据所有角色
func (s *service) SelectAll() (roles []*model.SystemRole, err error) {
	err = util.Database.Db.Model(&model.SystemRole{}).Find(&roles).Error
	return
}

// Select 查询角色
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.SystemRole, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemRole{})
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
	}
	if where["describe"] != "" {
		query.Where("`describe` like ?", "%"+where["describe"]+"%")
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
