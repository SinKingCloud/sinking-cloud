package config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// CountByKey 根据key查询用配置个数
func (s *service) CountByKey(key interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.SystemConfig{}).Where("`key` = ? ", key).Count(&total).Error
	return
}
