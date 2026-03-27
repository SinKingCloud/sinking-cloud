package config

import "service-user/app/model"

// Create 创建配置
func (r *Repository) Create(config *model.Config) error {
	return r.Database.Db.Create(config).Error
}
