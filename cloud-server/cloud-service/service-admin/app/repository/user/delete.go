package user

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteById 根据ID删除用户数据
func (r *Repository) DeleteById(id int64, tx ...*gorm.DB) error {
	if id <= 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return db.
			Where("id = ?", id).
			Delete(&model.User{}).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return execute(r.Database.Db)
}
