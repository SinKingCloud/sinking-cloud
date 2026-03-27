package config

import "service-admin/app/model"

// FindByIdentIdAndKey 根据identId和key查询配置
func (r *Repository) FindByIdentIdAndKey(identId, key string) ([]*model.Config, error) {
	var configs []*model.Config
	query := r.Database.Db.Model(&model.Config{})
	if identId != "" {
		query = query.Where("ident_id = ?", identId)
	}
	if key != "" {
		query = query.Where("key like ?", key+"%")
	}
	err := query.Find(&configs).Error
	return configs, err
}
