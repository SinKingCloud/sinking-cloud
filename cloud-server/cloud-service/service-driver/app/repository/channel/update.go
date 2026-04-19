package channel

import (
	"service-driver/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIds 根据 ids 更新
func (r *Repository) UpdateByIds(ids []int64, perm *UpdateChannel) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Type != nil {
		updates["type"] = perm.Type
	}
	if perm.Name != nil {
		updates["name"] = perm.Name
	}
	if perm.Config != nil {
		updates["config"] = perm.Config
	}
	if perm.Status != nil {
		updates["status"] = perm.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return tx.Model(&model.Channel{}).
				Where("id IN ?", batch.([]int64)).
				Updates(updates).Error
		})
	})
}
