package mind_share

import (
	"service-space/app/model"

	"gorm.io/gorm"
)

// DeleteByIdAndDataId 根据ID和导图数据ID删除分享信息
func (r *Repository) DeleteByIdAndDataId(id int64, dataId int64, tx ...*gorm.DB) error {
	if id <= 0 || dataId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Where("id = ? AND data_id = ?", id, dataId).Delete(&model.MindShare{}).Error
}

// DeleteByDataIds 根据导图数据ID集合删除分享信息
func (r *Repository) DeleteByDataIds(dataIds []int64, tx ...*gorm.DB) error {
	if len(dataIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(dataIds, 1000, func(batch interface{}) error {
			return db.Where("data_id IN ?", batch.([]int64)).Delete(&model.MindShare{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
