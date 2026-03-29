package user

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// FindById 根据id查询用户
func (r *Repository) FindById(id int64) (*model.User, error) {
	var user *model.User
	err := r.Database.Db.
		Where("id = ?", id).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user.Id <= 0 {
		return nil, nil
	}
	return user, nil
}

// FindByIdentifier 根据标识（手机号/邮箱/账户名）查询用户
func (r *Repository) FindByIdentifier(identifier string) (*model.User, error) {
	var user *model.User
	err := r.Database.Db.
		Where("phone = ? OR email = ? OR account = ?", identifier, identifier, identifier).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user.Id <= 0 {
		return nil, nil
	}
	return user, nil
}
