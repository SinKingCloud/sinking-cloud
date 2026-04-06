package role_permission

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByPermissionCodes 根据权限标识集合删除角色权限关联数据
func (r *Repository) DeleteByPermissionCodes(permissionCodes []string, tx ...*gorm.DB) error {
	if len(permissionCodes) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(permissionCodes, 1000, func(batch interface{}) error {
			return db.
				Where("permission_code IN ?", batch.([]string)).
				Delete(&model.RolePermission{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
