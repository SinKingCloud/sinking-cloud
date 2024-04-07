package system_menu

import (
	"encoding/json"
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

// Menu 页面输出结构
type Menu struct {
	Id    int    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Pid   int    `gorm:"column:pid" json:"pid"`
	Sort  int    `gorm:"column:sort" json:"sort"`
	Title string `gorm:"column:title" json:"title"`
	Key   string `gorm:"column:key" json:"key"`
	Path  string `gorm:"column:path" json:"path"`
	Icon  string `gorm:"column:icon" json:"icon"`
}

// SelectByIdsWithCache 根据菜单id查询角色菜单
func (s *Service) SelectByIdsWithCache(roleId int, menuIds string, refresh bool) (menus []*Menu, err error) {
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
func (s *Service) ClearRoleMenuCache(roleId int) bool {
	key := constant.CacheNameWithRoleMenu + strconv.Itoa(roleId)
	return util.Cache.Delete(key)
}
