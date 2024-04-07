package uin_info

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"time"
)

// UpdateByUin 通过uin更新
func (s *Service) UpdateByUin(uin string, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUinInfo{}).Where("`uin` = ? ", uin).Updates(data).Error
	return
}
