package constant

const (
	CacheNameWithUserConfig = "User_"    //用户应用配置
	CacheTimeWithUserConfig = 600 * 1000 //用户应用配置储存时间

	CacheNameWithSysConfig = "Sys"      //系统应用配置
	CacheTimeWithSysConfig = 600 * 1000 //用户应用配置储存时间

	CacheNameWithSession = "Session_"        //用户登陆会话
	CacheTimeWithSession = int64(600 * 1000) //用户登陆会话储存时间
)
