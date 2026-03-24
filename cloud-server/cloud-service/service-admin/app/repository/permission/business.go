package permission

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindByPath 根据路径查询权限
func (r *Repository) FindByPath(path string) (perm *model.Permission, err error) {
	err = r.Database.Db.Where("path = ?", path).First(&perm).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return perm, nil
}
