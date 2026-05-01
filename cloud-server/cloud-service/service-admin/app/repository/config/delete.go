package config

import (
	"errors"
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByIdentId 根据identId删除配置
func (r *Repository) DeleteByIdentId(identId string, tx ...*gorm.DB) error {
	if identId == "" {
		return errors.New("配置标识不能为空")
	}
	execute := func(db *gorm.DB) error {
		return db.Where("ident_id = ?", identId).Delete(&model.Config{}).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
