package system_access

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/constant"
	"sync"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

// FindByUriMethod 查询接口信息
func (*service) FindByUriMethod(url string, method string) (access *model.SystemAccess, err error) {
	err = util.Database.Db.Where("`url` = ? and `method`=?", url, method).First(&access).Error
	return
}

// FindByUriMethodWithCache 根据缓存查询接口信息
func (s *service) FindByUriMethodWithCache(url string, method string, refresh bool) (access *model.SystemAccess, err error) {
	key := constant.CacheNameWithSystemAccess + url + method
	if refresh {
		access, err = s.FindByUriMethod(url, method)
		if err != nil {
			return
		}
		data, err := json.Marshal(access)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemAccess)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &access)
		} else {
			access, err = s.FindByUriMethod(url, method)
			if err != nil {
				return
			}
			data, err := json.Marshal(access)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemAccess)
			}
		}
	}
	return
}
