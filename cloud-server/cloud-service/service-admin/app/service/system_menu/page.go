package system_menu

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/constant"
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

// Menu 页面输出结构
type Menu struct {
	Id    int    `json:"id"`
	Pid   int    `json:"pid"`
	Sort  int    `json:"sort"`
	Title string `json:"title"`
	Key   string `json:"key"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
}

// SelectByIdsWithCache 根据菜单id查询角色菜单
func (s *service) SelectByIdsWithCache(roleId int, menuIds string, refresh bool) (menus []*Menu, err error) {
	key := constant.CacheNameWithRoleMenu + strconv.Itoa(roleId)
	var ids []int
	err = json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return
	}
	if refresh {
		menus, err = s.SelectByIds(ids)
		if err != nil {
			return
		}
		data, err := json.Marshal(menus)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenu)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &menus)
		} else {
			menus, err = s.SelectByIds(ids)
			if err != nil {
				return
			}
			data, err := json.Marshal(menus)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenu)
			}
		}
	}
	return
}

// ClearRoleMenuCache 清理角色菜单缓存
func (s *service) ClearRoleMenuCache(roleId int) bool {
	key := constant.CacheNameWithRoleMenu + strconv.Itoa(roleId)
	return util.Cache.Delete(key)
}
