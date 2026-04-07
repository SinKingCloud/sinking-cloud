package user_role

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByUserId 根据用户ID删除用户角色关联
func (r *Repository) DeleteByUserId(userId int64, tx ...*gorm.DB) error {
	if userId <= 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return db.
			Where("user_id = ?", userId).
			Delete(&model.UserRole{}).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}

// DeleteByUserIds 根据用户ID集合删除用户角色关联
func (r *Repository) DeleteByUserIds(userIds []int64, tx ...*gorm.DB) error {
	if len(userIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(userIds, 1000, func(batch interface{}) error {
			return db.
				Where("user_id IN ?", batch.([]int64)).
				Delete(&model.UserRole{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
