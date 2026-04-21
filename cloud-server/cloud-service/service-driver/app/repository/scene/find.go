package scene

import (
	"errors"
	"service-driver/app/model"

	"gorm.io/gorm"
)

// FindById 根据ID查询记录
func (r *Repository) FindById(id int64) (*model.Scene, error) {
	var data *model.Scene
	err := r.Database.Db.Where("id = ?", id).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// ExistsByChannelId 判断channelId是否被使用
func (r *Repository) ExistsByChannelId(channelId int64) (bool, error) {
	var data *model.Scene
	err := r.Database.Db.
		Where("channel_id = ?", channelId).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return data != nil && data.Id > 0, nil
}
