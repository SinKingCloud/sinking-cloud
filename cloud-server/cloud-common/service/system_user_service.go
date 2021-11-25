package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"strconv"
)

// systemUserService 单例对象
type systemUserService struct {
}

// systemUserObj 单例对象
var systemUserObj *systemUserService

// GetSystemUserServiceIns 获取单例
func GetSystemUserServiceIns() *systemUserService {
	if systemUserObj == nil {
		systemUserObj = &systemUserService{}
	}
	return systemUserObj
}

// FindById 查询用户信息(根据用户ID)
func (systemUserService) FindById(id int) *model.SystemUser {
	var user *model.SystemUser
	util.Database.Db.Where("id = ?", id).First(&user)
	return user
}

// FindByIdWithCache 查询缓存用户信息(根据用户ID)
func (w systemUserService) FindByIdWithCache(id int, refresh bool) *model.SystemUser {
	var user *model.SystemUser
	key := util2.CacheNameWithWebUser + strconv.Itoa(id)
	if refresh {
		user = w.FindById(id)
		data, err := json.Marshal(user)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemUser)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &user)
		} else {
			user = w.FindById(id)
			data, err := json.Marshal(user)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemUser)
			}
		}
	}
	return user
}
