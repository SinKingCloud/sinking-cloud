package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"time"
)

// UpdateById 通过ID更新
func (s *Service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUin{}).Where("`id` = ? ", id).Updates(data).Error
	return
}

// UpdateByUin 通过uin更新
func (s *Service) UpdateByUin(uin string, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUin{}).Where("`uin` = ? ", uin).Updates(data).Error
	return
}

// UpdateByServerId 通过server_id更新
func (s *Service) UpdateByServerId(serverId int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUin{}).Where("`server_id` = ? ", serverId).Updates(data).Error
	return
}

// UpdateByIds 通过ID列表更新
func (s *Service) UpdateByIds(ids []int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.AppUin{}).Where("`id` in ? ", ids).Updates(data).Error
	return
}
