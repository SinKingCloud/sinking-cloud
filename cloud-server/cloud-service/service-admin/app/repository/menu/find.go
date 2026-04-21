package menu

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindByCode 根据菜单code查询信息
func (r *Repository) FindByCode(code string) (*model.Menu, error) {
	var menu *model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Where("code = ?", code).
		First(&menu).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if menu == nil || menu.Id <= 0 {
		return nil, nil
	}
	return menu, nil
}

// FindById 根据菜单id查询信息
func (r *Repository) FindById(id int64) (*model.Menu, error) {
	var menu *model.Menu
	err := r.Database.Db.
		Model(&model.Menu{}).
		Where("id = ?", id).
		First(&menu).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if menu == nil || menu.Id <= 0 {
		return nil, nil
	}
	return menu, nil
}
