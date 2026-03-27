package permission

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindByPath 根据路径查询权限
func (r *Repository) FindByPath(path string) (*model.Permission, error) {
	var perm *model.Permission
	err := r.Database.Db.Where("path = ?", path).First(&perm).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return perm, nil
}

// FindByCode 根据 code 查询权限
func (r *Repository) FindByCode(code string) (*model.Permission, error) {
	var perm *model.Permission
	err := r.Database.Db.Where("code = ?", code).First(&perm).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return perm, nil
}
