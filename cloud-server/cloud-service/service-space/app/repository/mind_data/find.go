package mind_data

import (
	"errors"
	"service-space/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindById 根据ID查询导图数据
func (r *Repository) FindById(id int64) (*model.MindData, error) {
	if id <= 0 {
		return nil, nil
	}
	var data *model.MindData
	err := r.Database.Db.Where("id = ?", id).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// FindByIdAndUserId 根据ID和用户ID查询导图数据
func (r *Repository) FindByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (*model.MindData, error) {
	if id <= 0 || userId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	lock := false
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
		lock = true
	}
	var data *model.MindData
	query := db.Where("id = ? AND user_id = ?", id, userId)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// ExistsByIdAndUserId 判断用户导图数据是否存在
func (r *Repository) ExistsByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (bool, error) {
	if id <= 0 || userId <= 0 {
		return false, nil
	}
	db := r.Database.Db
	lock := false
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
		lock = true
	}
	var data *MindData
	query := db.Model(&model.MindData{}).
		Select("id").
		Where("id = ? AND user_id = ?", id, userId)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return data != nil && data.Id > 0, nil
}
