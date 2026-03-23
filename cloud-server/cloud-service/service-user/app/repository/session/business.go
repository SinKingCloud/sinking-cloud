package session

import (
	"errors"
	"service-user/app/model"

	"gorm.io/gorm"
)

// Create 创建会话
func (r *Repository) Create(session *model.Session) error {
	return r.Database.Db.Create(session).Error
}

// DeleteByToken 通过token删除
func (r *Repository) DeleteByToken(token string) error {
	return r.Database.Db.Where("token = ?", token).Delete(&model.Session{}).Error
}

// FindByToken 根据token查询记录
func (r *Repository) FindByToken(token string) (*model.Session, error) {
	var session *model.Session
	err := r.Database.Db.
		Where("token = ?", token).
		First(&session).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return session, nil
}
