package session

import (
	"errors"
	"service-user/app/model"

	"gorm.io/gorm"
)

// FindByToken 根据token查询记录
func (r *Repository) FindByToken(token string) (*model.Session, error) {
	var session *model.Session
	err := r.Database.Db.
		Where("token = ?", token).
		First(&session).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if session == nil || session.Token == "" {
		return nil, nil
	}
	return session, nil
}
