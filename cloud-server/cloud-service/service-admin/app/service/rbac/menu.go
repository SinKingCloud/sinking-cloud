package rbac

import (
	"encoding/json"
	"service-admin/app/constant"
	"service-admin/app/enum/menu"
	"service-admin/app/model"
	"sort"
	"strconv"
)

// buildMenuTree 构建菜单树
func (s *service) buildMenuTree(menus []*model.Menu) []*MenuTree {
	menuMap := make(map[int64]*MenuTree, len(menus))
	var rootMenus []*MenuTree
	for _, item := range menus {
		menuMap[item.Id] = &MenuTree{
			Id:   item.Id,
			Code: item.Code,
			Name: item.Name,
			Type: item.Type,
			Icon: item.Icon,
			Sort: item.Sort,
		}
	}
	for _, item := range menus {
		current := menuMap[item.Id]
		if current == nil {
			continue
		}
		if item.ParentId == 0 {
			rootMenus = append(rootMenus, current)
			continue
		}
		if parentMenu, ok := menuMap[item.ParentId]; ok && parentMenu != nil {
			parentMenu.Children = append(parentMenu.Children, current)
			continue
		}
		rootMenus = append(rootMenus, current)
	}
	//菜单排序
	var sortMenus func(list []*MenuTree)
	sortMenus = func(list []*MenuTree) {
		if len(list) == 0 {
			return
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].Sort == list[j].Sort {
				return list[i].Id < list[j].Id
			}
			return list[i].Sort > list[j].Sort
		})
		for _, item := range list {
			if len(item.Children) > 0 {
				sortMenus(item.Children)
			}
		}
	}
	sortMenus(rootMenus)
	return rootMenus
}

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
