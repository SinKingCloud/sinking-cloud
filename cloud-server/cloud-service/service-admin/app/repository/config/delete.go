package config

import "service-admin/app/model"

// DeleteByIdentIdAndKey 根据identId和key删除配置
func (r *Repository) DeleteByIdentIdAndKey(identId, key string) error {
	query := r.Database.Db.Where("ident_id = ?", identId)
	if key != "" && key != "*" {
		query = query.Where("key like ?", key+"%")
	}
	return query.Delete(&model.Config{}).Error
}
