package mind_data

import (
	"service-space/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIdAndUserId 根据ID和用户ID更新导图数据
func (r *Repository) UpdateByIdAndUserId(id int64, userId int64, data *UpdateMindData, tx ...*gorm.DB) error {
	if id <= 0 || userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return r.update(db.Model(&model.MindData{}).
		Where("id = ? AND user_id = ?", id, userId), data)
}

// UpdateByTypeIdAndUserId 根据分类ID和用户ID批量更新导图数据
func (r *Repository) UpdateByTypeIdAndUserId(typeId int64, userId int64, data *UpdateMindData, tx ...*gorm.DB) error {
	if typeId <= 0 || userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return r.update(db.Model(&model.MindData{}).
		Where("type_id = ? AND user_id = ?", typeId, userId), data)
}

// update 更新导图数据
func (r *Repository) update(db *gorm.DB, data *UpdateMindData) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.TypeId != nil {
		updates["type_id"] = data.TypeId
	}
	if data.Title != nil {
		updates["title"] = data.Title
	}
	if data.Describe != nil {
		updates["describe"] = data.Describe
	}
	if data.Value != "" {
		updates["value"] = data.Value
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return db.Updates(updates).Error
}
