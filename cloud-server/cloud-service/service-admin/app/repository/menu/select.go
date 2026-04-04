package menu

import (
	"service-admin/app/model"
)

// SelectByIdsWithStatusAndType 根据菜单ID集合状态和类型查询
func (r *Repository) SelectByIdsWithStatusAndType(ids []int64, status int, types []int) ([]*model.Menu, error) {
	if len(ids) == 0 || len(types) == 0 {
		return nil, nil
	}
	var menuList []*model.Menu
	err := r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		batchIds := batch.([]int64)
		var menus []*model.Menu
		err := r.Database.Db.
			Model(&model.Menu{}).
			Where("id IN ? AND status = ? AND type IN ?", batchIds, status, types).
			Find(&menus).Error
		if err != nil {
			return err
		}
		menuList = append(menuList, menus...)
		return nil
	})
	return menuList, err
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
