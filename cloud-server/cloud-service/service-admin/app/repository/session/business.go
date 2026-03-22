package session

import "service-admin/app/model"

// Create 创建会话
func (r *Repository) Create(session *model.Session) error {
	return r.Database.Db.Create(session).Error
}
