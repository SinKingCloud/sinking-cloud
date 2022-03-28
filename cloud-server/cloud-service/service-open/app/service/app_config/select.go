package app_config

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
)

// SelectWithMap 查看map配置数据
func (*service) SelectWithMap(identId string, key string) map[string]string {
	var configs []*model.AppConfig
	temp := make(map[string]string)
	query := util.Database.Db
	if identId != "" {
		query.Where("`ident_id` = ?", identId)
	}
	if key != "" {
		query.Where("`key` = ?", key)
	}
	query.Find(&configs)
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}
