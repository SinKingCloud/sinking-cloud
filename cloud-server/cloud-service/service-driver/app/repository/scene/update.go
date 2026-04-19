package scene

import (
	"service-driver/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UpdateById 根据 id 更新
func (r *Repository) UpdateById(id int64, perm *UpdateScene) error {
	if perm == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if perm.Name != nil {
		updates["name"] = perm.Name
	}
	if perm.MaxSize != nil {
		updates["max_size"] = perm.MaxSize
	}
	if perm.AllowExt != nil {
		updates["allow_ext"] = perm.AllowExt
	}
	if perm.ChannelId != nil {
		updates["channel_id"] = perm.ChannelId
	}
	if perm.Remark != nil {
		updates["remark"] = perm.Remark
	}
	if perm.Status != nil {
		updates["status"] = perm.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Scene{}).Where("id = ?", id).Updates(updates).Error
}
