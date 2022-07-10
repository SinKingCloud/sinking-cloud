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

// ClearCache 清理缓存
func (systemConfigService) ClearCache() {
	util.Cache.Delete(util2.CacheNameWithSystemConfig)
	//查询所有group
	var configs []*model.SystemConfig
	e := util.Database.Db.Model(&model.SystemConfig{}).Where("`pid` = 0").Find(&configs).Error
	if e == nil {
		for _, v := range configs {
			util.Cache.Delete(util2.CacheNameWithSystemConfig + "_" + v.Key)
		}
	}
}

// SelectAllWithMap 获取map类型配置
func (systemConfigService) SelectAllWithMap() map[string]string {
	var configs []*model.SystemConfig
	temp := make(map[string]string)
	e := util.Database.Db.Find(&configs).Error
	if e == nil {
		for _, v := range configs {
			temp[v.Key] = v.Value
		}
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

// SelectByGroupWithMap 获取配置组获取
func (systemConfigService) SelectByGroupWithMap(group string) map[string]string {
	var configs []*model.SystemConfig
	table := &model.SystemConfig{}
	temp := make(map[string]string)
	e := util.Database.Db.Model(&table).Select("a.*").
		Joins("AS a LEFT JOIN `"+table.TableName()+"` as b ON a.`pid`=b.`id`").
		Where("b.`key` = ?", group).
		Scan(&configs).Error
	if e == nil {
		for _, v := range configs {
			temp[v.Key] = v.Value
		}
	}
	return temp
}

// GetGroupByCache 获取系统配置组
func (c systemConfigService) GetGroupByCache(group string, refresh bool) map[string]string {
	config := make(map[string]string)
	key := util2.CacheNameWithSystemConfig + "_" + group
	if refresh {
		config = c.SelectByGroupWithMap(group)
		data, err := json.Marshal(config)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemConfig)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &config)
		} else {
			config = c.SelectByGroupWithMap(group)
			data, err := json.Marshal(config)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemConfig)
			}
		}
	}
	return config
}
