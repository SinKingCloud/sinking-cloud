package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// Select 查询
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.AppUinServer, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AppUinServer{})
	if where["pid"] != "" {
		query.Where("`pid` = ?", where["pid"])
	}
	if where["key"] != "" {
		query.Where("`key` = ?", where["key"])
	}
	if where["name"] != "" {
		query.Where("`name` LIKE ?", "%"+where["name"]+"%")
	}
	if where["proxy"] != "" {
		query.Where("`proxy` LIKE ?", "%"+where["proxy"]+"%")
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["check_time_start"] != "" {
		query.Where("`check_time` >= ?", where["check_time_start"])
	}
	if where["check_time_end"] != "" {
		query.Where("`check_time` <= ?", where["check_time_end"])
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

// SelectServers 分组查询服务地址
func (s *Service) SelectServers(where map[string]string, page int, pageSize int) (list []*AppUinServer, total int64, err error) {
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
func (s *Service) SelectAllServer() (list []*model.AppUinServer, err error) {
	err = util.Database.Db.Model(&model.AppUinServer{}).Where("`status`=0 and `pid`>0").Find(&list).Error
	return
}
