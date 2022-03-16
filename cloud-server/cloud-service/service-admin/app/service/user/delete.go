package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
)

// DeleteByIds 通过ID列表删除
func (s *service) DeleteByIds(ids []int) (err error) {
	err = util.Database.Db.Where("`id` in ? and `id` != ?", ids, constant.AdminId).Delete(&model.SystemUser{}).Error
	return
}
