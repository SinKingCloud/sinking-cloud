package role_menu

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByMenuIds 根据菜单ID集合删除角色菜单关联数据
func (r *Repository) DeleteByMenuIds(menuIds []int64, tx ...*gorm.DB) error {
	if len(menuIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(menuIds, 1000, func(batch interface{}) error {
			return db.
				Where("menu_id IN ?", batch.([]int64)).
				Delete(&model.RoleMenu{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}

// DeleteByRoleIds 根据角色ID集合删除角色菜单关联数据
func (r *Repository) DeleteByRoleIds(roleIds []int64, tx ...*gorm.DB) error {
	if len(roleIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(roleIds, 1000, func(batch interface{}) error {
			return db.
				Where("role_id IN ?", batch.([]int64)).
				Delete(&model.RoleMenu{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
