package web_notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// Create 插入数据
func (s *Service) Create(data *model.WebNotice) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}
