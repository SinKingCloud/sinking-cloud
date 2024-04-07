package auth_version

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*AuthVersion, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.AuthVersion{})
	if where["auth_id"] != "" {
		query.Where("`auth_id` = ?", where["auth_id"])
	}
	if where["version"] != "" {
		query.Where("`version` = ?", where["version"])
	}
	if where["version_number"] != "" {
		query.Where("`version_number` = ?", where["version_number"])
	}
	if where["release_hash"] != "" {
		query.Where("`release_hash` = ?", where["release_hash"])
	}
	if where["update_hash"] != "" {
		query.Where("`update_hash` = ?", where["update_hash"])
	}
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["title"] != "" {
		query.Where("`title` like ?", "%"+where["title"]+"%")
	}
	if where["content"] != "" {
		query.Where("`content` like ?", "%"+where["content"]+"%")
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

// SelectUpgrade 查询大于版本
func (s *Service) SelectUpgrade(authId int, versionNumber int) (list []*AuthVersion, err error) {
	err = util.Database.Db.Model(&model.AuthVersion{}).
		Where("`auth_id` = ? and `version_number` > ? and `status` = ?", authId, versionNumber, Normal).
		Order("`version_number` desc").
		Find(&list).Error
	return
}
