package role_permission

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// Creates 批量创建角色权限关联
func (r *Repository) Creates(data []*model.RolePermission, tx ...*gorm.DB) error {
	if len(data) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(data, 1000, func(batch interface{}) error {
			return db.Create(batch.([]*model.RolePermission)).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
