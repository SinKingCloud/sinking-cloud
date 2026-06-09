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
