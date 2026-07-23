package mind_share

import (
	"service-space/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIdAndDataId 根据ID和导图数据ID更新分享信息
func (r *Repository) UpdateByIdAndDataId(id int64, dataId int64, data *UpdateMindShare, tx ...*gorm.DB) error {
	if id <= 0 || dataId <= 0 || data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if data.ExpireTime != nil {
		updates["expire_time"] = data.ExpireTime
	}
	if data.Remark != nil {
		updates["remark"] = data.Remark
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Model(&model.MindShare{}).
		Where("id = ? AND data_id = ?", id, dataId).
		Updates(updates).Error
}
