package menu

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	"strconv"
)

// SelectAllMenu 查询所有菜单
func SelectAllMenu() []*model.SystemMenu {
	var menus []*model.SystemMenu
	util.Database.Db.Model(&model.SystemMenu{}).Where("type=0 and status=0").Order("sort desc").Find(&menus)
	return menus
}

// SelectRoleControl 查询页面控件
func SelectRoleControl(menuIds []int, pid int) []*model.SystemMenu {
	var menus []*model.SystemMenu
	util.Database.Db.Model(&model.SystemMenu{}).Where("id in ? and pid = ? and type=1 and status=0", menuIds, pid).Order("sort desc").Find(&menus)
	return menus
}

// SelectRoleControlWithCache 查询页面控件(缓存)
func SelectRoleControlWithCache(roleId int, menuIds string, pid int, refresh bool) []*model.SystemMenu {
	var menus []*model.SystemMenu
	key := constant.CacheNameWithRoleMenuControl + strconv.Itoa(roleId) + "-" + strconv.Itoa(pid)
	var ids []int
	err := json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return nil
	}
	if refresh {
		menus = SelectRoleControl(ids, pid)
		data, err := json.Marshal(menus)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenuControl)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &menus)
		} else {
			menus = SelectRoleControl(ids, pid)
			data, err := json.Marshal(menus)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenuControl)
			}
		}
	}
	return menus
}

// SelectMenuByIds 根据id查询角色菜单
func SelectMenuByIds(menuIds []int) []*model.SystemMenu {
	var menus []*model.SystemMenu
	util.Database.Db.Model(&model.SystemMenu{}).Where("id in ? and type=0 and status=0", menuIds).Order("sort desc").Find(&menus)
	return menus
}

// SelectMenuByIdsWithCache 根据菜单id查询角色菜单
func SelectMenuByIdsWithCache(roleId int, menuIds string, refresh bool) []*model.SystemMenu {
	var menus []*model.SystemMenu
	key := constant.CacheNameWithRoleMenu + strconv.Itoa(roleId)
	var ids []int
	err := json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return nil
	}
	if refresh {
		menus = SelectMenuByIds(ids)
		data, err := json.Marshal(menus)
		if err == nil {
			util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenu)
		}
	} else {
		data, exists := util.Cache.Get(key)
		if exists {
			_ = json.Unmarshal([]byte(data), &menus)
		} else {
			menus = SelectMenuByIds(ids)
			data, err := json.Marshal(menus)
			if err == nil {
				util.Cache.SetWithExpire(key, string(data), constant.CacheTimeWithRoleMenu)
			}
		}
	}
	return menus
}
