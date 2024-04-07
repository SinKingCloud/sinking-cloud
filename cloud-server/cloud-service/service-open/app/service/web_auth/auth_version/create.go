package auth_version

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// Create 插入数据
func (s *Service) Create(data *model.AuthVersion) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}
