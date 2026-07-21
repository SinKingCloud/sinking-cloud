package mind_version

import (
	"errors"
	"service-space/app/model"

	"gorm.io/gorm"
)

// FindByIdAndDataId 根据ID和导图ID查询版本
func (r *Repository) FindByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) (*model.MindVersion, error) {
	if id <= 0 || dataId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var data *model.MindVersion
	err := db.Where("id = ? AND data_id = ?", id, dataId).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// FindLatestByDataId 查询导图数据最新版本
func (r *Repository) FindLatestByDataId(dataId int64, tx ...*gorm.DB) (*model.MindVersion, error) {
	if dataId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var data *model.MindVersion
	err := db.
		Where("data_id = ?", dataId).
		Order("id DESC").
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// FindLatestIdByDataId 查询导图数据最新版本ID
func (r *Repository) FindLatestIdByDataId(dataId int64) (int64, error) {
	if dataId <= 0 {
		return 0, nil
	}
	var data *MindVersion
	err := r.Database.Db.Model(&model.MindVersion{}).
		Select("id").
		Where("data_id = ?", dataId).
		Order("id DESC").
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if data == nil || data.Id <= 0 {
		return 0, nil
	}
	return data.Id, nil
}
