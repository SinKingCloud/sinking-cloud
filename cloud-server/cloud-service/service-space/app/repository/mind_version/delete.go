package mind_version

import (
	"service-space/app/model"

	"gorm.io/gorm"
)

// DeleteByDataIds 根据导图数据ID集合删除版本
func (r *Repository) DeleteByDataIds(dataIds []int64, tx ...*gorm.DB) error {
	if len(dataIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(dataIds, 1000, func(batch interface{}) error {
			return db.Where("data_id IN ?", batch.([]int64)).Delete(&model.MindVersion{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}

// DeleteOldByDataIdAndType 删除旧版本，仅保留最新 limit 条。
func (r *Repository) DeleteOldByDataIdAndType(dataId int64, versionType int, limit int, tx ...*gorm.DB) error {
	if dataId <= 0 || limit <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	where := map[string]interface{}{"data_id": dataId, "type": versionType}
	var total int64
	if err := db.Model(&model.MindVersion{}).Where(where).Count(&total).Error; err != nil || total <= int64(limit) {
		return err
	}
	var ids []int64
	if err := db.Model(&model.MindVersion{}).
		Where(where).
		Order("id ASC").
		Limit(int(total)-limit).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
		return db.Where("id IN ?", batch.([]int64)).Delete(&model.MindVersion{}).Error
	})
}
