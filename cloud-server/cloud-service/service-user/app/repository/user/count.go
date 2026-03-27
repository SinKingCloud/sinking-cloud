package user

import (
	"service-user/app/model"
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

// CountByEmail 根据邮箱统计用户数量
func (r *Repository) CountByEmail(email string) (int64, error) {
	var total int64
	err := r.Database.Db.Model(&model.User{}).
		Where("email = ?", email).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
