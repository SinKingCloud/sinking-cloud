package config

import "service-user/app/model"

// UpdateByUserIdAndKey 更新配置
func (r *Repository) UpdateByUserIdAndKey(userId int64, key string, value string) (int64, error) {
	result := r.Database.Db.Model(&model.Config{}).Where("user_id = ? and key = ?", userId, key).Update("value", value)
	return result.RowsAffected, result.Error
}
