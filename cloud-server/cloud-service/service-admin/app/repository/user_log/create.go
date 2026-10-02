package user_log

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// Create 创建用户操作日志
func (r *Repository) Create(data *model.UserLog, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(data).Error
}
