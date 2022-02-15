package service

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"strconv"
)

// webUserService 单例对象
type webUserService struct {
}

// userObj 单例对象
var webUserObj *webUserService

// GetWebUserServiceIns 获取单例
func GetWebUserServiceIns() *webUserService {
	if webUserObj == nil {
		webUserObj = &webUserService{}
	}
	return webUserObj
}

// ClearCache 清理缓存
func (webUserService) ClearCache(id int) {
	util.Cache.Delete(util2.CacheNameWithWebUser + strconv.Itoa(id))
}

// FindById 查询用户信息(根据用户ID)
func (webUserService) FindById(id int) *model.WebUser {
	var user *model.WebUser
	util.Database.Db.Where("id = ?", id).First(&user)
	return user
}

// FindByIdWithCache 查询缓存用户信息(根据用户ID)
func (w webUserService) FindByIdWithCache(id int, refresh bool) *model.WebUser {
	var user *model.WebUser
	key := util2.CacheNameWithWebUser + strconv.Itoa(id)
	if refresh {
		user = w.FindById(id)
		data, err := json.Marshal(user)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithWebUser)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &user)
		} else {
			user = w.FindById(id)
			data, err := json.Marshal(user)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), util2.CacheTimeWithWebUser)
			}
		}
	}
	return user
}
