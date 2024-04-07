package system_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// CountByIds 根据ids查询用户个数
func (s *Service) CountByIds(ids []int) (total int64, err error) {
	err = util.Database.Db.Model(&model.SystemUser{}).Where("`id` in ? ", ids).Count(&total).Error
	return
}

// CountByAccount 根据用户名查询用户个数
func (s *Service) CountByAccount(account interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.SystemUser{}).Where("`account` = ? ", account).Count(&total).Error
	return
}

// CountByPhone 根据手机号查询用户个数
func (s *Service) CountByPhone(data interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.SystemUser{}).Where("`phone` = ? ", data).Count(&total).Error
	return
}
