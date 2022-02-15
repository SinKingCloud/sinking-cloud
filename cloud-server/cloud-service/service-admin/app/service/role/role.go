package role

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	"strconv"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetRoleServiceIns 获取单例
func GetRoleServiceIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// FindById 查询角色信息(根据角色ID)
func (*service) FindById(id int) *model.SystemRole {
	var role *model.SystemRole
	util.Database.Db.Where("id = ? and status=0", id).First(&role)
	return role
}

// FindByIdWithCache 根据缓存查询角色信息(根据角色ID)
func (s *service) FindByIdWithCache(id int, refresh bool) *model.SystemRole {
	var role *model.SystemRole
	key := constant.CacheNameWithSystemRole + strconv.Itoa(id)
	if refresh {
		role = s.FindById(id)
		data, err := json.Marshal(role)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemRole)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &role)
		} else {
			role = s.FindById(id)
			data, err := json.Marshal(role)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithSystemRole)
			}
		}
	}
	return role
}
