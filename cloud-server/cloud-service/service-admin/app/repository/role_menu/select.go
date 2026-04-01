package role_menu

import "service-admin/app/model"

// SelectMenuIdsByRoleId 根据角色ID查询菜单ID集合
func (r *Repository) SelectMenuIdsByRoleId(roleId int64) ([]int64, error) {
	var menuIds []int64
	err := r.Database.Db.
		Model(&model.RoleMenu{}).
		Distinct("menu_id").
		Where("role_id = ?", roleId).
		Pluck("menu_id", &menuIds).Error
	if err != nil {
		return nil, err
	}
	return menuIds, nil
}
