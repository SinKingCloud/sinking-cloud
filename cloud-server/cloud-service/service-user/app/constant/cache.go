package constant

const (
	CacheNameWithUserConfig = "UserConfig_"     //用户应用配置
	CacheTimeWithUserConfig = int64(600 * 1000) //用户应用配置储存时间

	CacheNameWithSession = "Session_"        //用户登陆会话
	CacheTimeWithSession = int64(600 * 1000) //用户登陆会话储存时间
)
