package constant

const (
	CacheNameWithUserConfig = "SysUserConfig_" //用户应用配置
	CacheTimeWithUserConfig = 600 * 1000       //用户应用配置储存时间

	CacheNameWithSysConfig = "SysConfig_" //系统应用配置
	CacheTimeWithSysConfig = 600 * 1000   //用户应用配置储存时间

	CacheNameWithSession = "SysSession_"     //用户登陆会话
	CacheTimeWithSession = int64(600 * 1000) //用户登陆会话储存时间

	CacheNameWithUserInfo = "SysUserInfo_"    //用户信息
	CacheTimeWithUserInfo = int64(600 * 1000) //用户信息储存时间

	CacheNameWithPermission = "Permission_"     //权限信息缓存
	CacheTimeWithPermission = int64(600 * 1000) //权限信息缓存储存时间

	CacheNameWithUserRoles = "UserRoles_"      //用户角色
	CacheTimeWithUserRoles = int64(600 * 1000) //用户角色信息储存时间

	CacheNameWithRolePermissions = "RolePermissions_" //角色权限
	CacheTimeWithRolePermissions = int64(600 * 1000)  //角色权限储存时间
)
