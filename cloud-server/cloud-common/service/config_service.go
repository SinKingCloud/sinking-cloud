package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// configService config_service
type configService struct {
}

// configObj config对象
var configObj *configService

// GetConfigServiceIns 获取单例
func GetConfigServiceIns() *configService {
	if configObj == nil {
		configObj = &configService{}
	}
	return configObj
}

// SelectAllConfigWithMap 获取map类型配置
func (configService) SelectAllConfigWithMap() map[string]string {
	var configs []*model.SystemConfig
	temp := make(map[string]string)
	util.Database.Db.Find(&configs)
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}

// GetConfigByCache 获取系统配置
func (c configService) GetConfigByCache(refresh bool) map[string]string {
	config := make(map[string]string)
	if refresh {
		config = c.SelectAllConfigWithMap()
		data, err := json.Marshal(config)
		if err == nil {
			util.Cache.SetWithExpire(util2.CacheNameWithSystemConfig, string(data), util2.CacheTimeWithSystemConfig)
		}
	} else {
		data, exists := util.Cache.Get(util2.CacheNameWithSystemConfig)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = c.SelectAllConfigWithMap()
			data, err := json.Marshal(config)
			if err == nil {
				util.Cache.SetWithExpire(util2.CacheNameWithSystemConfig, string(data), util2.CacheTimeWithSystemConfig)
			}
		}
	}
	return config
}
