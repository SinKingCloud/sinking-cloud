package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"sync"
)

// WebUserConfService config_service
type WebUserConfService struct {
}

// webUserConfigObj systemUserConfigObj对象
var (
	webUserConfigObj     *WebUserConfService
	onceWebUserConfigObj sync.Once
)

// GetWebUserConfSrvIns 获取单例
func GetWebUserConfSrvIns() *WebUserConfService {
	onceWebUserConfigObj.Do(func() {
		webUserConfigObj = &WebUserConfService{}
	})
	return webUserConfigObj
}

// GetConfigs 获取用户所有配置项
func (s *WebUserConfService) GetConfigs(userId int, prefix string) map[string]string {
	data := make(map[string]string)
	var conf []*model.WebUserConfig
	err := util.Database.Db.Where("`user_id` = ? AND `key` LIKE ?", userId, prefix+"%").Find(&conf).Error
	if err == nil {
		for _, v := range conf {
			data[v.Key] = v.Value
		}
	}
	return data
}

// SetConfig 设置配置
func (s *WebUserConfService) SetConfig(userId int, key string, value string) bool {
	var total int64
	err := util.Database.Db.Where("`user_id` = ? and `key` = ?", userId, key).Updates(&model.WebUserConfig{
		Value: value,
	}).Count(&total).Error
	if err != nil || total <= 0 {
		err = util.Database.Db.Create(&model.WebUserConfig{
			UserId: userId,
			Key:    key,
			Value:  value,
		}).Error
	}
	return err == nil
}

// SetConfigs 设置配置
func (s *WebUserConfService) SetConfigs(userId int, data map[string]string) {
	for k, v := range data {
		s.SetConfig(userId, k, v)
	}
}

// FindByKey 通过key查询配置
func (s *WebUserConfService) FindByKey(key string) (config *model.WebUserConfig, err error) {
	err = util.Database.Db.Where("`key` = ?", key).Find(&config).Error
	return
}

// FindByValue 通过value查询配置
func (s *WebUserConfService) FindByValue(value string) (config *model.WebUserConfig, err error) {
	err = util.Database.Db.Where("`value` = ?", value).Find(&config).Error
	return
}
