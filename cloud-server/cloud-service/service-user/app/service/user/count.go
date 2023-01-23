package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// CountByAccount 根据用户名查询用户个数
func (s *service) CountByAccount(account interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.WebUser{}).Where("`account` = ? ", account).Count(&total).Error
	return
}

// CountByPhone 根据手机号查询用户个数
func (s *service) CountByPhone(phone interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.WebUser{}).Where("`phone` = ? ", phone).Count(&total).Error
	return
}
