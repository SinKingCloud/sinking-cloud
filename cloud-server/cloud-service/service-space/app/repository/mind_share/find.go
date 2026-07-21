package mind_share

import (
	"errors"
	"service-space/app/model"

	"gorm.io/gorm"
)

// FindByIdAndDataId 根据ID和导图数据ID查询分享信息
func (r *Repository) FindByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) (*model.MindShare, error) {
	if id <= 0 || dataId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var data *model.MindShare
	err := db.Where("id = ? AND data_id = ?", id, dataId).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// FindByCode 根据分享码查询分享信息
func (r *Repository) FindByCode(code string, tx ...*gorm.DB) (*model.MindShare, error) {
	if code == "" {
		return nil, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var data *model.MindShare
	err := db.Where("code = ?", code).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}
