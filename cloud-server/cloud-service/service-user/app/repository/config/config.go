package config

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Repository 配置仓储
type Repository struct {
	Database *database.Database
}

// FindByUserIdAndKey 根据userId和key查询
func (r *Repository) FindByUserIdAndKey(userId int64, key string) ([]*model.Config, error) {
	var configs []*model.Config
	query := r.Database.Db.Model(&model.Config{})
	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}
	if key != "" {
		query = query.Where("key like ?", key+"%")
	}
	err := query.Find(&configs).Error
	return configs, err
}

// UpdateByUserIdAndKey 更新配置
func (r *Repository) UpdateByUserIdAndKey(userId int64, key string, value string) (int64, error) {
	result := r.Database.Db.Model(&model.Config{}).Where("user_id = ? and key = ?", userId, key).Update("value", value)
	return result.RowsAffected, result.Error
}

// Create 创建配置
func (r *Repository) Create(config *model.Config) error {
	return r.Database.Db.Create(config).Error
}
