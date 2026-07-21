package mind_type

import (
	"errors"
	"service-space/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindByIdAndUserId 根据ID和用户ID查询导图分类
func (r *Repository) FindByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) (*model.MindType, error) {
	if id <= 0 || userId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	lock := false
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
		lock = true
	}
	var data *model.MindType
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

// FindByNameAndUserId 根据名称和用户ID查询导图分类
func (r *Repository) FindByNameAndUserId(name string, userId int64, tx ...*gorm.DB) (*model.MindType, error) {
	if name == "" || userId <= 0 {
		return nil, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var data *model.MindType
	err := db.Where("name = ? AND user_id = ?", name, userId).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}
