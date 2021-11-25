package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
)

// systemConfigService config_service
type systemConfigService struct {
}

// systemConfigObj systemConfigObj对象
var systemConfigObj *systemConfigService

// GetSystemConfigServiceIns 获取单例
func GetSystemConfigServiceIns() *systemConfigService {
	if systemConfigObj == nil {
		systemConfigObj = &systemConfigService{}
	}
	return systemConfigObj
}

// SelectAllWithMap 获取map类型配置
func (systemConfigService) SelectAllWithMap() map[string]string {
	var configs []*model.SystemConfig
	temp := make(map[string]string)
	util.Database.Db.Find(&configs)
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}

// GetByCache 获取系统配置
func (c systemConfigService) GetByCache(refresh bool) map[string]string {
	config := make(map[string]string)
	if refresh {
		config = c.SelectAllWithMap()
		data, err := json.Marshal(config)
		if err == nil {
			util.Cache.SetWithExpire(util2.CacheNameWithSystemConfig, string(data), util2.CacheTimeWithSystemConfig)
		}
	} else {
		data, exists := util.Cache.Get(util2.CacheNameWithSystemConfig)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = c.SelectAllWithMap()
			data, err := json.Marshal(config)
			if err == nil {
				util.Cache.SetWithExpire(util2.CacheNameWithSystemConfig, string(data), util2.CacheTimeWithSystemConfig)
			}
		}
	}
	return config
}
