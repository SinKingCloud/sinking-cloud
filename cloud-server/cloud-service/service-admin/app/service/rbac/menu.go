package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/enum/menu"
	"service-admin/app/model"
	"strconv"
)

// getMenuByCode 根据code查询菜单
func (s *service) getMenuByCode(code string, refresh bool) (menu *model.Menu, err error) {
	cacheKey := constant.CacheNameWithMenus + code
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &menu); err == nil {
				return menu, nil
			}
		}
	}
	menu, err = s.repositoryMenu.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(menu); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithMenus)
	}
	return menu, nil
}

// getMenuByCode 根据code查询菜单
func (s *service) getMenusByParentId(parentId int64, refresh bool) (menus []*model.Menu, err error) {
	cacheKey := constant.CacheNameWithMenus + "_list_" + strconv.FormatInt(parentId, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			if err = json.Unmarshal([]byte(data), &menus); err == nil {
				return menus, nil
			}
		}
	}
	menus, err = s.repositoryMenu.SelectByParentIdWithStatusAndType(parentId, int(menu.StatusEnable), int(menu.TypeControl))
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(menus); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithMenus)
	}
	return menus, nil
}
