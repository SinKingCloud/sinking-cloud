package balance

import (
	"errors"
	"service-pay/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindOrCreateByUserId 查询记录(不存在则创建)
func (r *Repository) FindOrCreateByUserId(userId int64, tx ...*gorm.DB) (*model.Balance, error) {
	var balance *model.Balance
	execute := func(db *gorm.DB) error {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(&model.Balance{
			UserId: userId,
		}).Error
		if err != nil {
			return err
		}
		query := db.Where("user_id = ?", userId)
		if len(tx) > 0 && tx[0] != nil {
			query = query.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return query.First(&balance).Error
	}
	var err error
	if len(tx) > 0 && tx[0] != nil {
		err = execute(tx[0])
	} else {
		err = r.Database.Transaction(execute)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if balance == nil || balance.Id <= 0 {
		return nil, nil
	}
	return balance, nil
}

// SelectByUserIds 通过用户ID批量查询余额
func (r *Repository) SelectByUserIds(userIds []int64) ([]*model.Balance, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	var data []*model.Balance
	err := r.Database.BatchExecute(userIds, 1000, func(batch interface{}) error {
		var list []*model.Balance
		err := r.Database.Db.
			Model(&model.Balance{}).
			Where("user_id IN ?", batch.([]int64)).
			Find(&list).Error
		if err != nil {
			return err
		}
		data = append(data, list...)
		return nil
	})
	return data, err
}
