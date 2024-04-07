package cookie

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/model"
	"strconv"
)

// FindByUinDomain 通过uin domain查询
func (s *Service) FindByUinDomain(uin string, domain string) (data *model.AppUinCookie, err error) {
	err = util.Database.Db.Where("`uin` = ? AND `domain` = ?", uin, domain).First(&data).Error
	return
}

// FindByUinDomainWithCache 通过uin domain查询(缓存)
func (s *Service) FindByUinDomainWithCache(uin string, domain string, refresh bool) (data *model.AppUinCookie, err error) {
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
func (s *Service) ClearByUinDomain(uin string, domain string) bool {
	key := constant.CacheNameWithUinCookie + uin + domain
	return util.Cache.Delete(key)
}

// FindByStatusDomain 通过status domain查询
func (s *Service) FindByStatusDomain(status int, domain string) (data *model.AppUinCookie, err error) {
	err = util.Database.Db.Where("`status` = ? AND `domain` = ?", status, domain).Order("`check_time` desc").First(&data).Error
	return
}

// FindByStatusDomainWithCache 通过status domain查询(缓存)
func (s *Service) FindByStatusDomainWithCache(status int, domain string, refresh bool) (data *model.AppUinCookie, err error) {
	key := constant.CacheNameWithUinCookie + strconv.Itoa(status) + domain
	if refresh {
		data, err = s.FindByStatusDomain(status, domain)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), 600000)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindByStatusDomain(status, domain)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), 600000)
			}
		}
	}
	return
}

// FindByStatusDomainOrderByIdDesc 通过status domain查询
func (s *Service) FindByStatusDomainOrderByIdDesc(status int, domain string) (data *model.AppUinCookie, err error) {
	err = util.Database.Db.Where("`status` = ? AND `domain` = ?", status, domain).Order("`id` desc").First(&data).Error
	return
}

// FindByStatusDomainOrderByIdWithCache 通过status domain查询(缓存)
func (s *Service) FindByStatusDomainOrderByIdWithCache(status int, domain string, refresh bool) (data *model.AppUinCookie, err error) {
	key := constant.CacheNameWithUinCookie + strconv.Itoa(status) + domain
	if refresh {
		data, err = s.FindByStatusDomain(status, domain)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), 600000)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindByStatusDomain(status, domain)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), 600000)
			}
		}
	}
	return
}

// ClearByStatusDomain 清除缓存
func (s *Service) ClearByStatusDomain(status int, domain string) bool {
	key := constant.CacheNameWithUinCookie + strconv.Itoa(status) + domain
	return util.Cache.Delete(key)
}
