package role_menu

import "service-admin/app/model"

import "gorm.io/gorm"

// DeleteByMenuIds 根据菜单ID集合删除角色菜单关联数据
func (r *Repository) DeleteByMenuIds(menuIds []int64) error {
	if len(menuIds) == 0 {
		return nil
	}
	return r.Database.Db.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(menuIds, 1000, func(batch interface{}) error {
			return tx.
				Where("menu_id IN ?", batch.([]int64)).
				Delete(&model.RoleMenu{}).Error
		})
	})
}
