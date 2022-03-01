package access

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// FindByUriMethod 查询接口信息
func (*service) FindByUriMethod(url string, method string) *model.SystemAccess {
	var access *model.SystemAccess
	util.Database.Db.Where("url = ? and method=? and status=0", url, method).First(&access)
	return access
}

// FindByUriMethodWithCache 根据缓存查询接口信息
func (s *service) FindByUriMethodWithCache(url string, method string, refresh bool) *model.SystemAccess {
	var access *model.SystemAccess
	key := constant.CacheNameWithSystemAccess + url + method
	if refresh {
		access = s.FindByUriMethod(url, method)
		data, err := json.Marshal(access)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemAccess)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &access)
		} else {
			access = s.FindByUriMethod(url, method)
			data, err := json.Marshal(access)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemAccess)
			}
		}
	}
	return access
}
