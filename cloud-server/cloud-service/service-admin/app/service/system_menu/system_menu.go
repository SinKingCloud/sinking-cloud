package system_menu

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	"strconv"
)

// FindByKeyWithCache 通过key查找页面数据（缓存）
func (s *service) FindByKeyWithCache(key string, refresh bool) (menu *Menu, err error) {
	cacheKey := constant.CacheNameRoleMenu + key
	if refresh {
		menu, err = s.FindByKey(key)
		if err != nil {
			return
		}
		data, err := json.Marshal(menu)
		if err == nil {
			util.Cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeRoleMenu)
		}
	} else {
		data, exists := util.Cache.Get(cacheKey)
		if exists {
			_ = json.Unmarshal([]byte(data), &menu)
		} else {
			menu, err = s.FindByKey(key)
			if err != nil {
				return
			}
			data, err := json.Marshal(menu)
			if err == nil {
				util.Cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeRoleMenu)
			}
		}
	}
	return
}

// SelectControlWithCache 查询页面控件(缓存)
func (s *service) SelectControlWithCache(roleId int, menuIds string, pid int, refresh bool) (menus []*Menu, err error) {
	key := constant.CacheNameWithRoleMenuControl + strconv.Itoa(roleId) + "-" + strconv.Itoa(pid)
	var ids []int
	err = json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return
	}
	if refresh {
		menus, err = s.SelectControl(ids, pid)
		if err != nil {
			return
		}
		data, err := json.Marshal(menus)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenuControl)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &menus)
		} else {
			menus, err = s.SelectControl(ids, pid)
			if err != nil {
				return
			}
			data, err := json.Marshal(menus)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenuControl)
			}
		}
	}
	return
}

// GetAccessIdsByMenuIds 通过菜单获取关联接口权限
func (s *service) GetAccessIdsByMenuIds(ids []int) (accessIds []int, err error) {
	menus, err := s.SelectFullByIds(ids)
	if err != nil {
		return
	}
	for _, v := range menus {
		var tempAccessIds []int
		tempErr := json.Unmarshal([]byte(v.Accesses), &tempAccessIds)
		if tempErr != nil {
			continue
		}
		accessIds = append(accessIds, tempAccessIds...)
	}
	return
}
