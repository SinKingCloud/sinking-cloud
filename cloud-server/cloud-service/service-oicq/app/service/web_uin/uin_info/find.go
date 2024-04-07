package uin_info

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
)

// FindByUin 通过ID查询
func (s *Service) FindByUin(uin string) (data *model.AppUinInfo, err error) {
	err = util.Database.Db.Where("`uin` = ? ", uin).First(&data).Error
	return
}
