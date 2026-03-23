package session

import "service-admin/app/model"

// Create 创建会话
func (r *Repository) Create(session *model.Session) error {
	return r.Database.Db.Create(session).Error
}

// DeleteByToken 通过token删除
func (r *Repository) DeleteByToken(token string) error {
	return r.Database.Db.Where("token = ?", token).Delete(&model.Session{}).Error
}
