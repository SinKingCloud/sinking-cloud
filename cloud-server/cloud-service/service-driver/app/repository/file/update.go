package file

import (
	"service-driver/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIds 根据 ids 更新数据
func (r *Repository) UpdateByIds(ids []int64, perm *UpdateFile) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Status != nil {
		updates["status"] = perm.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return tx.Model(&model.File{}).
				Where("id IN ?", batch.([]int64)).
				Updates(updates).Error
		})
	})
}
