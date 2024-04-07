package system_file

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Select 查询角色
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.SystemFile, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.SystemFile{})
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["hash"] != "" {
		query.Where("`hash` = ?", where["hash"])
	}
	if where["ext"] != "" {
		query.Where("`ext` = ?", where["ext"])
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

// SelectPathByIds 通过id查询文件path
func (s *Service) SelectPathByIds(ids []int) []string {
	var list []string
	util.Database.Db.Model(&model.SystemFile{}).Where("`id` in ? ", ids).Select("`path`").Scan(&list)
	return list
}
