package role

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindById 根据ID查询记录
func (r *Repository) FindById(id int64) (*model.Role, error) {
	var role *model.Role
	err := r.Database.Db.Where("id = ?", id).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if role.Id <= 0 {
		return nil, nil
	}
	return role, nil
}
