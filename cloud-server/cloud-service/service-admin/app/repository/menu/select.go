package menu

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// SelectByIds 根据菜单ID集合查询
func (r *Repository) SelectByIds(ids []int64) ([]*model.Menu, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var menuList []*model.Menu
	err := r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		batchIds := batch.([]int64)
		var menus []*model.Menu
		err := r.Database.Db.
			Model(&model.Menu{}).
			Where("id IN ?", batchIds).
			Find(&menus).Error
		if err != nil {
			return err
		}
		menuList = append(menuList, menus...)
		return nil
	})
	return menuList, err
}

// SelectAll 查询全部菜单数据
func (r *Repository) SelectAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// SelectWithStatus 根据状态查询全部菜单数据
func (r *Repository) SelectWithStatus(status int) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Where("status = ?", status).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

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

// SelectWithType 根据菜单父级ID和类型查询子项
func (r *Repository) SelectWithType(types []int) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Where("type IN ?", types).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// Select 查询菜单数据
func (r *Repository) Select(where *SelectMenu, queryPage *page.Query) (*page.Result[*model.Menu], error) {
	query := r.Database.Db.Model(&model.Menu{})
	if where != nil {
		if where.ParentId != "" {
			query = query.Where("parent_id = ?", where.ParentId)
		}
		if where.Code != "" {
			query = query.Where("code LIKE ?", "%"+where.Code+"%")
		}
		if where.Name != "" {
			query = query.Where("name LIKE ?", "%"+where.Name+"%")
		}
		if where.Icon != "" {
			query = query.Where("icon LIKE ?", "%"+where.Icon+"%")
		}
		if where.Type != "" {
			query = query.Where("type = ?", where.Type)
		}
		if where.Status != "" {
			query = query.Where("status = ?", where.Status)
		}
		if where.CreateTimeStart != "" {
			query = query.Where("create_time >= ?", where.CreateTimeStart)
		}
		if where.CreateTimeEnd != "" {
			query = query.Where("create_time <= ?", where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != "" {
			query = query.Where("update_time >= ?", where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != "" {
			query = query.Where("update_time <= ?", where.UpdateTimeEnd)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
