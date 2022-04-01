package system_config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Create 插入数据
func (s *service) Create(data *model.SystemConfig) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}
