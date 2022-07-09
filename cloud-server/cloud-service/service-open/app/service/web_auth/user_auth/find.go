package user_auth

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"strconv"
)

// Find 查询
func (s *service) Find(where map[string]string) (data *model.UserAuth, err error) {
	query := util.Database.Db.Model(&model.UserAuth{})
	if where["user_id"] != "" {
		query.Where("`user_id` = ?", where["user_id"])
	}
	if where["auth_id"] != "" {
		query.Where("`auth_id` = ?", where["auth_id"])
	}
	if where["domain"] != "" {
		query.Where("`domain` = ?", where["domain"])
	}
	if where["code"] != "" {
		query.Where("`code` = ?", where["code"])
	}
	if where["ip"] != "" {
		query.Where("`ip` = ?", where["ip"])
	}
	err = query.First(&data).Error
	return
}

// FindByAuthIdDomain 通过程序ID以及域名查询
func (s *service) FindByAuthIdDomain(authId int, domain string) (data *model.UserAuth, err error) {
	err = util.Database.Db.Where("`auth_id` = ? and `domain` = ?", authId, domain).First(&data).Error
	return
}

// FindById 通过ID查询
func (s *service) FindById(id int) (data *model.UserAuth, err error) {
	err = util.Database.Db.Where("`id` = ? ", id).First(&data).Error
	return
}

// FindByIdWithCache 通过ID查询
func (s *service) FindByIdWithCache(id int, refresh bool) (data *model.UserAuth, err error) {
	key := constant.CacheNameWithUserAuth + strconv.Itoa(id)
	if refresh {
		data, err = s.FindById(id)
		if err != nil {
			return
		}
		str, err := json.Marshal(data)
		if err == nil {
			util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithUserAuth)
		}
	} else {
		str, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(str), &data)
		} else {
			data, err = s.FindById(id)
			if err != nil {
				return
			}
			str, err := json.Marshal(data)
			if err == nil {
				util.Cache.SetWithExpire(key, string(str), constant.CacheTimeWithUserAuth)
			}
		}
	}
	return
}

// Clear 清除缓存
func (s *service) Clear(id int) bool {
	key := constant.CacheNameWithAuthVersion + strconv.Itoa(id)
	return util.Cache.Delete(key)
}
