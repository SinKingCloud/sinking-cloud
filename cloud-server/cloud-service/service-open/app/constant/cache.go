package constant

const (
	CacheNameWithAppKey = "AppKey_"  //应用app储存名称
	CacheTimeWithAppKey = 600 * 1000 //应用app储存时间

	CacheNameWithOauth = "Oauth_"   //oauth储存名称
	CacheTimeWithOauth = 600 * 1000 //oauth储存时间

	CacheNameWithAuth = "Auth_"    //授权auth储存名称
	CacheTimeWithAuth = 600 * 1000 //授权auth储存时间

	CacheNameWithAuthVersion = "AuthVersion_" //授权auth_version储存名称
	CacheTimeWithAuthVersion = 600 * 1000     //授权auth_version储存时间

	CacheNameWithUserAuth = "UserAuth_" //授权user_auth储存名称
	CacheTimeWithUserAuth = 600 * 1000  //授权user_auth储存时间

	CacheNameWithAppConfig = "AppConfig_" //应用配置
	CacheTimeWithAppConfig = 600 * 1000   //应用配置储存时间

	CacheNameWithUserAppConfig = "UserAppConfig_" //用户应用配置
	CacheTimeWithUserAppConfig = 600 * 1000       //用户应用配置储存时间

	CacheNameWithUinServer = "UinServer_" //挂机服务器信息
	CacheTimeWithUinServer = 600 * 1000   //挂机服务器信息储存时间

	CacheNameWithRbac = "Rbac_"    //rbac权限信息缓存
	CacheTimeWithRbac = 600 * 1000 //默认rbac权限信息缓存时间
)
