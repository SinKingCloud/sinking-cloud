package config

import "service-admin/app/model"

// UpdateByIdentIdAndKey 更新配置
func (r *Repository) UpdateByIdentIdAndKey(identId, key, value string) (int64, error) {
	result := r.Database.Db.Model(&model.Config{}).Where("ident_id = ? and key = ?", identId, key).Update("value", value)
	return result.RowsAffected, result.Error
}
