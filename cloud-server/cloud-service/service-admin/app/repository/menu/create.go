package menu

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// Create 创建数据
func (r *Repository) Create(data *model.Menu, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Create(data).Error
	}
	return r.Database.Db.Create(data).Error
}
