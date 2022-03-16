package constant

const (
	CacheTime = 600 //默认缓存时间

	CacheNameWithSystemRole = "SystemRole_"    //权限角色信息缓存
	CacheTimeWithSystemRole = CacheTime * 1000 //默认权限角色缓存时间

	CacheNameWithSystemAccess = "SystemAccess_"  //权限接口信息缓存
	CacheTimeWithSystemAccess = CacheTime * 1000 //默认权限接口信息缓存时间

	CacheNameWithRbac = "Rbac_"          //rbac权限信息缓存
	CacheTimeWithRbac = CacheTime * 1000 //默认rbac权限信息缓存时间

	CacheNameRoleMenu = "Menu_"          //权限菜单缓存
	CacheTimeRoleMenu = CacheTime * 1000 //默认权限菜单缓存时间

	CacheNameWithRoleMenu = "RoleMenu_"      //权限菜单缓存
	CacheTimeWithRoleMenu = CacheTime * 1000 //默认权限菜单缓存时间

	CacheNameWithRoleMenuControl = "RoleMenuControl_" //权限菜单控件缓存
	CacheTimeWithRoleMenuControl = CacheTime * 1000   //默认权限菜单控件缓存时间
)
