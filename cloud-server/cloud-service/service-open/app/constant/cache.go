package constant

const (
	CacheNameWithAppKey = "AppKey_"  //应用app储存名称
	CacheTimeWithAppKey = 600 * 1000 //应用app储存时间

	CacheNameWithOauth = "Oauth_"   //oauth储存名称
	CacheTimeWithOauth = 600 * 1000 //oauth储存时间

	CacheNameWithAppConfig = "AppConfig_" //应用配置
	CacheTimeWithAppConfig = 600 * 1000   //应用配置储存时间

	CacheNameWithUserAppConfig = "UserAppConfig_" //用户应用配置
	CacheTimeWithUserAppConfig = 600 * 1000       //用户应用配置储存时间

	CacheNameWithUinServer = "UinServer_" //挂机服务器信息
	CacheTimeWithUinServer = 600 * 1000   //挂机服务器信息储存时间
)
