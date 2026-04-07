package session

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByToken 通过token删除
func (r *Repository) DeleteByToken(token string) error {
	return r.Database.Db.Where("token = ?", token).Delete(&model.Session{}).Error
}

// DeleteByUserId 根据用户ID删除会话
func (r *Repository) DeleteByUserId(userId int64, tx ...*gorm.DB) error {
	execute := func(db *gorm.DB) error {
		return db.
			Where("user_id = ?", userId).
			Delete(&model.Session{}).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
