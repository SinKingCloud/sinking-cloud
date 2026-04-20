package scene

import (
	"service-driver/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UpdateById 根据 id 更新
func (r *Repository) UpdateById(id int64, data *UpdateScene) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Name != nil {
		updates["name"] = data.Name
	}
	if data.MaxSize != nil {
		updates["max_size"] = data.MaxSize
	}
	if data.AllowExt != nil {
		updates["allow_ext"] = data.AllowExt
	}
	if data.ChannelId != nil {
		updates["channel_id"] = data.ChannelId
	}
	if data.Remark != nil {
		updates["remark"] = data.Remark
	}
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Scene{}).Where("id = ?", id).Updates(updates).Error
}
