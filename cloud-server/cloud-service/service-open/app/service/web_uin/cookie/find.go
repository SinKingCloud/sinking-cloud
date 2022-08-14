package cookie

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
)

// FindByUinDomain 通过uin domain查询
func (s *service) FindByUinDomain(uin string, domain string) (data *model.AppUinCookie, err error) {
	err = util.Database.Db.Where("`uin` = ? AND `domain` = ?", uin, domain).First(&data).Error
	return
}

// FindByUinDomainWithCache 通过uin domain查询(缓存)
func (s *service) FindByUinDomainWithCache(uin string, domain string, refresh bool) (data *model.AppUinCookie, err error) {
	key := constant.CacheNameWithUinCookie + uin + domain
	if refresh {
		data, err = s.FindByUinDomain(uin, domain)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithUinCookie)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindByUinDomain(uin, domain)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithUinCookie)
			}
		}
	}
	return
}

// ClearByUinDomain 清除缓存
func (s *service) ClearByUinDomain(uin string, domain string) bool {
	key := constant.CacheNameWithUinCookie + uin + domain
	return util.Cache.Delete(key)
}
