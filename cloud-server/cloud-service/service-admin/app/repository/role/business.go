package role

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindById 根据ID查询记录
func (r *Repository) FindById(id int64) (role *model.Role, err error) {
	err = r.Database.Db.Where("id = ?", id).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return role, nil
}
