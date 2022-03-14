package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// CountByIds 根据ids查询角色菜单
func (s *service) CountByIds(ids []int) (total int64, err error) {
	err = util.Database.Db.Model(&model.SystemUser{}).Where("`id` in ? ", ids).Count(&total).Error
	return
}
