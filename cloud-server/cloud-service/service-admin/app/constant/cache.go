package constant

const (
	CacheNameWithUserConfig = "SysUserConfig_" //用户应用配置
	CacheTimeWithUserConfig = 600 * 1000       //用户应用配置储存时间

	CacheNameWithSysConfig = "SysConfig_" //系统配置
	CacheTimeWithSysConfig = 600 * 1000   //系统配置储存时间

	CacheNameWithSession = "SysSession_"     //用户登陆会话
	CacheTimeWithSession = int64(600 * 1000) //用户登陆会话储存时间

	CacheNameWithUserInfo = "SysUserInfo_"    //用户信息
	CacheTimeWithUserInfo = int64(600 * 1000) //用户信息储存时间

	CacheNameWithUserVerify = "SysUserVerify_" //安全验证

	CacheNameWithPermission = "Permission_"     //权限信息缓存
	CacheTimeWithPermission = int64(600 * 1000) //权限信息缓存储存时间

	CacheNameWithUserRoles = "UserRoles_"      //用户角色
	CacheTimeWithUserRoles = int64(600 * 1000) //用户角色信息储存时间

	CacheNameWithRolePermissions = "RolePermissions_" //角色权限
	CacheTimeWithRolePermissions = int64(600 * 1000)  //角色权限储存时间

	CacheNameWithRoleMenus = "RoleMenus_"      //角色菜单
	CacheTimeWithRoleMenus = int64(600 * 1000) //角色菜单储存时间

	CacheNameWithMenus = "Menus_"          //菜单信息
	CacheTimeWithMenus = int64(600 * 1000) //菜单信息储存时间

	CacheNameWithEnum = "Enum_"           //枚举缓存
	CacheTimeWithEnum = int64(600 * 1000) //枚举缓存储存时间

	CacheNameWithDashboardServiceCount = "DashboardServiceCount" //首页服务数量
	CacheTimeWithDashboardServiceCount = int64(1 * 60 * 1000)    //首页服务数量储存时间
)
