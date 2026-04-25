package channel

import (
	"service-pay/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UpdateById 根据 id 更新
func (r *Repository) UpdateById(id int64, perm *UpdateChannel) error {
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
	return r.Database.Db.Model(&model.Channel{}).Where("id = ?", id).Updates(updates).Error
}
