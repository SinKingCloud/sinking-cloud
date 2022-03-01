package menu

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

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
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

// SelectAllMenu 查询所有菜单
func (s *service) SelectAllMenu() (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("type=0 and status=0").Order("sort desc").Find(&menus).Error
	return
}

// SelectRoleControl 查询页面控件
func (s *service) SelectRoleControl(menuIds []int, pid int) (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("id in ? and pid = ? and type=1 and status=0", menuIds, pid).Order("sort desc").Find(&menus).Error
	return
}

// SelectRoleControlWithCache 查询页面控件(缓存)
func (s *service) SelectRoleControlWithCache(roleId int, menuIds string, pid int, refresh bool) (menus []*Menu, err error) {
	key := constant.CacheNameWithRoleMenuControl + strconv.Itoa(roleId) + "-" + strconv.Itoa(pid)
	var ids []int
	err = json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return
	}
	if refresh {
		menus, err = s.SelectRoleControl(ids, pid)
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
			menus, err = s.SelectRoleControl(ids, pid)
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

// SelectMenuByIds 根据id查询角色菜单
func (s *service) SelectMenuByIds(menuIds []int) (menus []*Menu, err error) {
	err = util.Database.Db.Model(&model.SystemMenu{}).Where("id in ? and type=0 and status=0", menuIds).Order("sort desc").Find(&menus).Error
	return
}

// SelectMenuByIdsWithCache 根据菜单id查询角色菜单
func (s *service) SelectMenuByIdsWithCache(roleId int, menuIds string, refresh bool) (menus []*Menu, err error) {
	key := constant.CacheNameWithRoleMenu + strconv.Itoa(roleId)
	var ids []int
	err = json.Unmarshal([]byte(menuIds), &ids)
	if err != nil {
		return
	}
	if refresh {
		menus, err = s.SelectMenuByIds(ids)
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
			menus, err = s.SelectMenuByIds(ids)
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
