package user

import (
	"errors"
	"service-user/app/model"

	"gorm.io/gorm"
)

// CountByPhone 根据手机号统计用户数量
func (r *Repository) CountByPhone(phone string) (int64, error) {
	var total int64
	err := r.Database.Db.Model(&model.User{}).
		Where("phone = ?", phone).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
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
	return user, nil
}

// Create 创建用户
func (r *Repository) Create(user *model.User) error {
	return r.Database.Db.Create(user).Error
}
