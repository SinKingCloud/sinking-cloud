package constant

const (
	CacheNameWithUserInfo = "WebUserInfo_"    //用户信息
	CacheTimeWithUserInfo = int64(600 * 1000) //用户信息储存时间

	CacheNameWithUserVerify = "UserVerify_" //安全验证

	CacheNameWithUserConfig = "UserConfig_"     //用户应用配置
	CacheTimeWithUserConfig = int64(600 * 1000) //用户应用配置储存时间

	CacheNameWithSession = "Session_"        //用户登陆会话
	CacheTimeWithSession = int64(600 * 1000) //用户登陆会话储存时间

	CacheNameWithRealAuthInfo = "RealAuthInfo_"   //实名认证信息
	CacheTimeWithRealAuthInfo = int64(600 * 1000) //实名认证信息储存时间

	CacheNameWithRealAuthProcess = "RealAuthProcess_"    //实名认证流程
	CacheTimeWithRealAuthProcess = int64(30 * 60 * 1000) //实名认证流程储存时间
)
