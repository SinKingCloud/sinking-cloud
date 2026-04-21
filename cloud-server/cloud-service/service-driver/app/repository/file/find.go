package file

import (
	"errors"
	"service-driver/app/model"

	"gorm.io/gorm"
)

// ExistsByChannelId 判断channelId是否被使用
func (r *Repository) ExistsByChannelId(channelId int64) (bool, error) {
	var data *model.File
	err := r.Database.Db.
		Where("channel_id = ?", channelId).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return data != nil && data.Id > 0, nil
}

// ExistsBySceneId 判断sceneId是否被使用
func (r *Repository) ExistsBySceneId(sceneId int64) (bool, error) {
	var data *model.File
	err := r.Database.Db.
		Where("scene_id = ?", sceneId).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return data != nil && data.Id > 0, nil
}
