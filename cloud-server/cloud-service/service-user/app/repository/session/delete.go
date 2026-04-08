package session

import (
	"service-user/app/model"
)

// DeleteByToken 通过token删除
func (r *Repository) DeleteByToken(token string) error {
	return r.Database.Db.Where("token = ?", token).Delete(&model.Session{}).Error
}

// DeleteByUserIdAndTokens 根据用户ID和token集合删除会话
func (r *Repository) DeleteByUserIdAndTokens(userId int64, tokens []string) error {
	if userId <= 0 || len(tokens) == 0 {
		return nil
	}
	return r.Database.BatchExecute(tokens, 1000, func(batch interface{}) error {
		return r.Database.Db.
			Where("user_id = ?", userId).
			Where("token IN ?", batch.([]string)).
			Delete(&model.Session{}).Error
	})
}
