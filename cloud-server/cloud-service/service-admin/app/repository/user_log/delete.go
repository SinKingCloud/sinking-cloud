package user_log

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByUserId 根据用户ID删除操作日志
func (r *Repository) DeleteByUserId(userId int64, tx ...*gorm.DB) error {
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Where("user_id = ?", userId).Delete(&model.UserLog{}).Error
}
