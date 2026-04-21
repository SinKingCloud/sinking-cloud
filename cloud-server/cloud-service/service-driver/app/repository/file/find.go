package file

import (
	"errors"
	"service-driver/app/model"

	"gorm.io/gorm"
)

// FindByHashWithSceneId 根据hash和场景查询文件。
func (r *Repository) FindByHashWithSceneId(hash string, sceneId int64) (*model.File, error) {
	var data *model.File
	err := r.Database.Db.
		Where("hash = ? AND scene_id = ?", hash, sceneId).
		First(&data).Error
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
