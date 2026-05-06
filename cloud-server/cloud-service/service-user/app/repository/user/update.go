package user

import (
	"service-user/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UpdateById 根据ID更新
func (r *Repository) UpdateById(id int64, data *UpdateUser) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Account != nil {
		updates["account"] = data.Account
	}
	if data.Phone != nil {
		updates["phone"] = data.Phone
	}
	if data.Email != nil {
		updates["email"] = data.Email
	}
	if data.Password != nil {
		updates["password"] = data.Password
	}
	if data.Avatar != nil {
		updates["avatar"] = data.Avatar
	}
	if data.NickName != nil {
		updates["nick_name"] = data.NickName
	}
	if data.RealAuthType != nil {
		updates["real_auth_type"] = data.RealAuthType
	}
	if data.RealAuthStatus != nil {
		updates["real_auth_status"] = data.RealAuthStatus
	}
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
