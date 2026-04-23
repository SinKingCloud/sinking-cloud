package balance

import (
	"errors"
	"service-pay/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindOrCreateByUserId 查询记录(不存在则创建)
func (r *Repository) FindOrCreateByUserId(userId int64) (*model.Balance, error) {
	var balance *model.Balance
	err := r.Database.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(&model.Balance{
			UserId: userId,
		}).Error
		if err != nil {
			return err
		}
		return tx.Where("user_id = ?", userId).First(&balance).Error
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if balance == nil || balance.Id <= 0 {
		return nil, nil
	}
	return balance, nil
}
