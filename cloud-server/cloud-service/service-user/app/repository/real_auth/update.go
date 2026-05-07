package real_auth

import (
	"service-user/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateById 根据ID更新实名记录。
func (r *Repository) UpdateById(id int64, data *UpdateRealAuth, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Type != nil {
		updates["type"] = data.Type
	}
	if data.CompanyName != nil {
		updates["company_name"] = data.CompanyName
	}
	if data.CompanyNo != nil {
		updates["company_no"] = data.CompanyNo
	}
	if data.CardName != nil {
		updates["card_name"] = data.CardName
	}
	if data.CardNo != nil {
		updates["card_no"] = data.CardNo
	}
	if data.Files != nil {
		updates["files"] = data.Files
	}
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if data.AuthTime != nil {
		updates["auth_time"] = data.AuthTime
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Model(&model.RealAuth{}).Where("id = ?", id).Updates(updates).Error
}
