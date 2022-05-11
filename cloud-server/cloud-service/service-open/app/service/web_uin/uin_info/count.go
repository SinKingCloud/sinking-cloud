package uin_info

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// CountByUin 统计
func (s *service) CountByUin(uin string) (total int64) {
	e := util.Database.Db.Model(&model.AppUinInfo{}).Where("`uin` = ? ", uin).Count(&total).Error
	if e != nil {
		total = 0
	}
	return
}
