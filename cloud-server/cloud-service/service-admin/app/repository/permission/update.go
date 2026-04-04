package permission

import (
	"service-admin/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByCode 根据 code 更新权限
func (r *Repository) UpdateByCode(code string, perm *UpdatePermission) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Name != "" {
		updates["name"] = perm.Name
	}
	if perm.Path != "" {
		updates["path"] = perm.Path
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Permission{}).Where("code = ?", code).Updates(updates).Error
}

// UpdateByIds 根据 ids 更新权限
func (r *Repository) UpdateByIds(ids []int64, perm *UpdatePermission) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Name != "" {
		updates["name"] = perm.Name
	}
	if perm.Path != "" {
		updates["path"] = perm.Path
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return tx.Model(&model.Permission{}).
				Where("id IN ?", batch.([]int64)).
				Updates(updates).Error
		})
	})
}
