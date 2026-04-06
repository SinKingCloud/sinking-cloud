package menu

import (
	"time"

	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateById 根据ID更新数据
func (r *Repository) UpdateById(id int64, data *UpdateMenu, tx ...*gorm.DB) error {
	if id <= 0 || data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.ParentId != nil {
		updates["parent_id"] = data.ParentId
	}
	if data.Name != nil {
		updates["name"] = data.Name
	}
	if data.Icon != nil {
		updates["icon"] = data.Icon
	}
	if data.Sort != nil {
		updates["sort"] = data.Sort
	}
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	execute := func(db *gorm.DB) error {
		return db.
			Model(&model.Menu{}).
			Where("id = ?", id).
			Updates(updates).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return execute(r.Database.Db)
}
