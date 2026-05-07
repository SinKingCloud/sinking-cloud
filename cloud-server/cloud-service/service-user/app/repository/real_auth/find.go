package real_auth

import (
	"errors"
	"service-user/app/enum/real_auth_status"
	"service-user/app/enum/real_auth_type"
	"service-user/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FindOrCreateByUserId 根据用户ID查询实名记录，不存在则创建。
func (r *Repository) FindOrCreateByUserId(userId int64, tx ...*gorm.DB) (*model.RealAuth, error) {
	var info *model.RealAuth
	execute := func(db *gorm.DB) error {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(&model.RealAuth{
			UserId: userId,
			Type:   real_auth_type.UnAuth,
			Status: real_auth_status.UnAuth,
		}).Error; err != nil {
			return err
		}
		query := db.Where("user_id = ?", userId)
		if len(tx) > 0 && tx[0] != nil {
			query = query.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return query.First(&info).Error
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
	if info == nil || info.Id <= 0 {
		return nil, nil
	}
	return info, nil
}
