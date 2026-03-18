package config

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository Repository实现
type Repository struct {
	Database *database.Database
}

// FindByIdentIdAndKey 根据identId和key查询配置
func (r *Repository) FindByIdentIdAndKey(identId, key string) ([]*model.Config, error) {
	var configs []*model.Config
	query := r.Database.Db.Model(&model.Config{})
	if identId != "" {
		query = query.Where("`ident_id` = ?", identId)
	}
	if key != "" {
		query = query.Where("`key` like ?", key+"%")
	}
	err := query.Find(&configs).Error
	return configs, err
}

// UpdateByIdentIdAndKey 更新配置
func (r *Repository) UpdateByIdentIdAndKey(identId, key, value string) (int64, error) {
	result := r.Database.Db.Model(&model.Config{}).Where("`ident_id` = ? and `key` = ?", identId, key).Update("value", value)
	return result.RowsAffected, result.Error
}

// Create 创建配置
func (r *Repository) Create(config *model.Config) error {
	return r.Database.Db.Create(config).Error
}
