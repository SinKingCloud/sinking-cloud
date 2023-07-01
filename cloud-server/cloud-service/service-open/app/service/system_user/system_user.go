package system_user

import (
	"encoding/json"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"strconv"
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

// FindById 查询用户信息(根据用户ID)
func (*service) FindById(id int) (user *model.SystemUser) {
	user = rpc.SystemUser.FindById(id)
	return
}

// FindByIdWithCache 获取用户信息(通过ID缓存)
func (s *service) FindByIdWithCache(id int, refresh bool) (user *model.SystemUser, err error) {
	key := util2.CacheNameWithSystemUser + strconv.Itoa(id)
	if refresh {
		user = s.FindById(id)
		if user != nil {
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
			user = s.FindById(id)
			if user != nil {
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

// GetCurrentUserInfo  获取登录用户最新信息
func (s *service) GetCurrentUserInfo(id int, refresh bool) *model.SystemUser {
	user, _ := s.FindByIdWithCache(id, refresh)
	return user
}
