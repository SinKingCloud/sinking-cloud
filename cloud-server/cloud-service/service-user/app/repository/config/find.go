package config

import "service-user/app/model"

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
