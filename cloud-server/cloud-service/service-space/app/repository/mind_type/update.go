package mind_type

import (
	"service-space/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIdAndUserId 根据ID和用户ID更新导图分类
func (r *Repository) UpdateByIdAndUserId(id int64, userId int64, name string, tx ...*gorm.DB) error {
	if id <= 0 || userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Model(&model.MindType{}).
		Where("id = ? AND user_id = ?", id, userId).
		Updates(map[string]interface{}{
			"name":        name,
			"update_time": str.DateTime(time.Now()),
		}).Error
}
