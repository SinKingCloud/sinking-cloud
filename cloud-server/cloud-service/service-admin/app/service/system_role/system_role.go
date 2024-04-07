package system_role

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/constant"
	"strconv"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// FindByIdWithCache 根据缓存查询角色信息(根据角色ID)
func (s *Service) FindByIdWithCache(id int, refresh bool) (role *model.SystemRole, err error) {
	key := constant.CacheNameWithSystemRole + strconv.Itoa(id)
	if refresh {
		role, err = s.FindById(id)
		if err != nil {
			return
		}
		data, err := json.Marshal(role)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemRole)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &role)
		} else {
			role, err = s.FindById(id)
			if err != nil {
				return
			}
			data, err := json.Marshal(role)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemRole)
			}
		}
	}
	return
}

// ClearRoleCache 清理角色缓存(根据角色ID)
func (s *Service) ClearRoleCache(id int) bool {
	key := constant.CacheNameWithSystemRole + strconv.Itoa(id)
	return util.Cache.Delete(key)
}
