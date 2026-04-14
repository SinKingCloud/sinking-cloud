package task

import (
	"service-cron/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UpdateById 根据ID更新
func (r *Repository) UpdateById(id int64, data *UpdateTask) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Type != nil {
		updates["type"] = data.Type
	}
	if data.Name != nil {
		updates["phone"] = data.Name
	}
	if data.Spec != nil {
		updates["email"] = data.Spec
	}
	if data.Script != nil {
		updates["password"] = data.Script
	}
	if data.Concurrent != nil {
		updates["avatar"] = data.Concurrent
	}
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Task{}).Where("id = ?", id).Updates(updates).Error
}
