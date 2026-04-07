package role

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// Create 创建角色
func (r *Repository) Create(data *model.Role, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Create(data).Error
	}
	return r.Database.Db.Create(data).Error
}
