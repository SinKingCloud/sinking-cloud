package app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Create 插入数据
func (s *service) Create(data *model.WebUserApp) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}
