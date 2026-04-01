package menu

import (
	"service-admin/app/model"
)

// SelectByIdsWithStatusAndType 根据菜单ID集合状态和类型查询
func (r *Repository) SelectByIdsWithStatusAndType(ids []int64, status int, types []int) ([]*model.Menu, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	total := len(ids)
	var menuList []*model.Menu
	const batchSize = 1000
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batchIds := ids[i:end]
		var menus []*model.Menu
		err := r.Database.Db.
			Model(&model.Menu{}).
			Where("id IN ? AND status = ? AND type IN ?", batchIds, status, types).
			Find(&menus).Error
		if err != nil {
			return nil, err
		}
		menuList = append(menuList, menus...)
	}
	return menuList, nil
}

// SelectByParentIdWithStatusAndType 根据菜单父级ID和类型查询子项
func (r *Repository) SelectByParentIdWithStatusAndType(parentId int64, status int, types int) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Where("parent_id = ? AND status = ? AND type = ?", parentId, status, types).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
