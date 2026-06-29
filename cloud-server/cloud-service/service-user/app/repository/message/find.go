package message

import (
	"errors"
	"service-user/app/model"

	"gorm.io/gorm"
)

// FindByIdAndUserId 根据消息ID和用户ID查询消息详情
func (r *Repository) FindByIdAndUserId(id int64, userId int64) (*model.Message, error) {
	var info *model.Message
	err := r.Database.Db.
		Where("id = ? AND user_id = ?", id, userId).
		First(&info).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if info == nil || info.Id <= 0 {
		return nil, nil
	}
	return info, nil
}
