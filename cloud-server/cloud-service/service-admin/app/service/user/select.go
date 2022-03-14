package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// SelectByIds 根据id查询用户
func (s *service) SelectByIds(menuIds []int) (users []*model.SystemUser, err error) {
	err = util.Database.Db.Model(&model.SystemUser{}).Where("`id` in ? ", menuIds).Find(&users).Error
	return
}
