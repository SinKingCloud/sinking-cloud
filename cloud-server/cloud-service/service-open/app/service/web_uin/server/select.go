package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 查询
func (s *service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.AppUinServer, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppUinServer{})
	if where["key"] != "" {
		query.Where("`key` = ?", where["key"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["name"] != "" {
		query.Where("`name` like ?", "%"+where["name"]+"%")
	}
	if where["addr"] != "" {
		query.Where("`addr` like ?", "%"+where["addr"]+"%")
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}

// SelectServers 分组查询服务地址
func (s *service) SelectServers(where map[string]string, page int, pageSize int) (list []*AppUinServer, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppUinServer{}).Select("`key`,`name`,`status`")
	if where["pid"] != "" {
		query.Where("`pid` = ?", where["pid"])
	}
	if where["key"] != "" {
		query.Where("`key` = ?", where["key"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	err = query.Group("`key`").Group("`name`").Group("`status`").Count(&total).Limit(pageSize).Offset(offset).Scan(&list).Error
	return
}

// SelectAllServer 查询所有代理服务器
func (s *service) SelectAllServer() (list []*model.AppUinServer, err error) {
	err = util.Database.Db.Model(&model.AppUinServer{}).Where("`status`=0 and `pid`>0").Find(&list).Error
	return
}
