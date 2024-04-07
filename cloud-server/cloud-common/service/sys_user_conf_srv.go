package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"sync"
)

// SysUserConfService config_service
type SysUserConfService struct {
}

// systemConfigObj systemUserConfigObj对象
var (
	systemUserConfigObj     *SysUserConfService
	onceSystemUserConfigObj sync.Once
)

// GetSysUserConfSrvIns 获取单例
func GetSysUserConfSrvIns() *SysUserConfService {
	onceSystemUserConfigObj.Do(func() {
		systemUserConfigObj = &SysUserConfService{}
	})
	return systemUserConfigObj
}

// GetConfigs 获取用户所有配置项
func (s *SysUserConfService) GetConfigs(userId int, prefix string) map[string]string {
	data := make(map[string]string)
	var conf []*model.SystemUserConfig
	err := util.Database.Db.Where("`user_id` = ? AND `key` LIKE ?", userId, prefix+"%").Find(&conf).Error
	if err == nil {
		for _, v := range conf {
			data[v.Key] = v.Value
		}
	}
	return data
}

// SetConfig 设置配置
func (s *SysUserConfService) SetConfig(userId int, key string, value string) bool {
	var total int64
	err := util.Database.Db.Where("`user_id` = ? and `key` = ?", userId, key).Updates(&model.SystemUserConfig{
		Value: value,
	}).Count(&total).Error
	if err != nil || total <= 0 {
		err = util.Database.Db.Create(&model.SystemUserConfig{
			UserId: userId,
			Key:    key,
			Value:  value,
		}).Error
	}
	return err == nil
}

// SetConfigs 设置配置
func (s *SysUserConfService) SetConfigs(userId int, data map[string]string) {
	for k, v := range data {
		s.SetConfig(userId, k, v)
	}
}

// FindByKey 通过key查询配置
func (s *SysUserConfService) FindByKey(key string) (config *model.SystemUserConfig, err error) {
	err = util.Database.Db.Where("`key` = ?", key).Find(&config).Error
	return
}

// FindByValue 通过value查询配置
func (s *SysUserConfService) FindByValue(value string) (config *model.SystemUserConfig, err error) {
	err = util.Database.Db.Where("`value` = ?", value).Find(&config).Error
	return
}
