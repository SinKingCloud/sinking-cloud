package user_plugin

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"strconv"
)

// FindByUserIdAndPluginId 通过用户ID和插件ID查询
func (s *Service) FindByUserIdAndPluginId(userId int, pluginId int) (data *model.UserPlugin, err error) {
	err = util.Database.Db.Where("`user_id` = ? AND `user_id` = ?", userId, pluginId).First(&data).Error
	return
}

// FindByUserIdAndPluginIdWithCache 通过用户ID和插件ID查询
func (s *Service) FindByUserIdAndPluginIdWithCache(userId int, pluginId int, refresh bool) (data *model.UserPlugin, err error) {
	key := constant.CacheNameWithPlugin + strconv.Itoa(userId) + "_" + strconv.Itoa(pluginId)
	if refresh {
		data, err = s.FindByUserIdAndPluginId(userId, pluginId)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithPlugin)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindByUserIdAndPluginId(userId, pluginId)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithPlugin)
			}
		}
	}
	return
}

// ClearByUserIdAndPluginId 清除缓存
func (s *Service) ClearByUserIdAndPluginId(userId int, pluginId int) bool {
	key := constant.CacheNameWithPlugin + strconv.Itoa(userId) + "_" + strconv.Itoa(pluginId)
	return util.Cache.Delete(key)
}
