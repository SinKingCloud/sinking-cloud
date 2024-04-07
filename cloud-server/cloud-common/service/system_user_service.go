package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"strconv"
	"sync"
)

// SystemUserService 单例对象
type SystemUserService struct {
}

// systemUserObj 单例对象
var (
	systemUserObj     *SystemUserService
	onceSystemUserObj sync.Once
)

// GetSystemUserServiceIns 获取单例
func GetSystemUserServiceIns() *SystemUserService {
	onceSystemUserObj.Do(func() {
		systemUserObj = &SystemUserService{}
	})
	return systemUserObj
}

// ClearCache 清理缓存
func (SystemUserService) ClearCache(id int) {
	util.Cache.Delete(util2.CacheNameWithSystemUser + strconv.Itoa(id))
}

// FindById 查询用户信息(根据用户ID)
func (SystemUserService) FindById(id int) (user *model.SystemUser, err error) {
	err = util.Database.Db.Where("id = ?", id).First(&user).Error
	return
}

// FindByIdWithCache 查询缓存用户信息(根据用户ID) refresh为true强制刷新
func (w SystemUserService) FindByIdWithCache(id int, refresh bool) (user *model.SystemUser, err error) {
	key := util2.CacheNameWithSystemUser + strconv.Itoa(id)
	if refresh {
		user, err = w.FindById(id)
		if err != nil {
			return
		}
		data, err := json.Marshal(user)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemUser)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			err = json.Unmarshal([]byte(data), &user)
		} else {
			user, err = w.FindById(id)
			if err != nil {
				return
			}
			data, err := json.Marshal(user)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithSystemUser)
			}
		}
	}
	return
}
