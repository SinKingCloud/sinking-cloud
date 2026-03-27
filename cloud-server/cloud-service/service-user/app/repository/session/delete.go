package session

import "service-user/app/model"

// DeleteByToken 通过token删除
func (r *Repository) DeleteByToken(token string) error {
	return r.Database.Db.Where("token = ?", token).Delete(&model.Session{}).Error
}
