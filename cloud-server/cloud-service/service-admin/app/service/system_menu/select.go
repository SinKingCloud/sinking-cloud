package system_menu

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// SelectControl 查询页面控件
func (s *service) SelectControl(menuIds []int, pid int) (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("`id` in ? and `pid` = ? and `type`=1 and `status`=0", menuIds, pid).Order("sort desc").Find(&menus).Error
	return
}

// SelectByIds 根据id查询菜单
func (s *service) SelectByIds(menuIds []int) (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("`id` in ? and `type`=0 and `status`=0", menuIds).Order("`sort` desc").Find(&menus).Error
	return
}

// SelectFullByIds 根据id查询菜单全部信息
func (s *service) SelectFullByIds(menuIds []int) (menus []*model.SystemMenu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("`id` in ?", menuIds).Find(&menus).Error
	return
}

// SelectAll 查询所有菜单
func (s *service) SelectAll() (menus []*model.SystemMenu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Order("`sort` desc").Find(&menus).Error
	return
}

// SelectByPid 通过pid查找数据
func (s *service) SelectByPid(pid int) (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("`pid` = ?", pid).Find(&menus).Error
	return
}

// Select 查询菜单
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.SystemMenu, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemMenu{})
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["title"] != "" {
		query.Where("`title` like ?", "%"+where["title"]+"%")
	}
	if where["key"] != "" {
		query.Where("`key` like ?", "%"+where["key"]+"%")
	}
	if where["path"] != "" {
		query.Where("`path` like ?", "%"+where["path"]+"%")
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
