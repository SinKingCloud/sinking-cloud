package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"time"
)

// UpdateById 通过ID更新
func (s *service) UpdateById(id int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.SystemConfig{}).Where("`id` = ? ", id).Updates(data).Error
	return
}

// UpdateByKey 通过KEY更新
func (s *service) UpdateByKey(key string, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.SystemConfig{}).Where("`key` = ? ", key).Updates(data).Error
	return
}

// UpdateByIds 通过ID列表更新
func (s *service) UpdateByIds(ids []int, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.SystemConfig{}).Where("`id` in ? ", ids).Updates(data).Error
	return
}
