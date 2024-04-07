package plugin

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// FindByUrl 通过URL标识查询
func (s *Service) FindByUrl(url string) (data *model.AppPlugin, err error) {
	err = util.Database.Db.Where("`filter_url` = ? ", url).First(&data).Error
	return
}

// FindByUrlWithCache 通过URL查询
func (s *Service) FindByUrlWithCache(url string, refresh bool) (data *model.AppPlugin, err error) {
	key := constant.CacheNameWithPlugin + url
	if refresh {
		data, err = s.FindByUrl(url)
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
			data, err = s.FindByUrl(url)
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

// ClearByUrl 通过URL标识清除缓存
func (s *Service) ClearByUrl(url string) bool {
	key := constant.CacheNameWithPlugin + url
	return util.Cache.Delete(key)
}
