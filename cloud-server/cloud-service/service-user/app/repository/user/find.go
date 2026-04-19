package user

import (
	"errors"
	"service-user/app/model"

	"gorm.io/gorm"
)

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

// ExistsByPhone 判断手机号是否存在
func (r *Repository) ExistsByPhone(phone string) (bool, error) {
	var user *model.User
	err := r.Database.Db.
		Where("phone = ?", phone).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return user.Id > 0, nil
}

// ExistsByEmail 判断邮箱是否存在
func (r *Repository) ExistsByEmail(email string) (bool, error) {
	var user *model.User
	err := r.Database.Db.
		Where("email = ?", email).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return user.Id > 0, nil
}

// ExistsByAccount 判断账号是否存在
func (r *Repository) ExistsByAccount(account string) (bool, error) {
	var user *model.User
	err := r.Database.Db.
		Where("account = ?", account).
		First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return user.Id > 0, nil
}
