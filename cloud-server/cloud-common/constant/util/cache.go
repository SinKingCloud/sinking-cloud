package util

const (
	CacheNameWithSystemConfig = "SystemConfig" //系统配置缓存
	CacheTimeWithSystemConfig = 86400 * 1000   //系统配置缓存时间

	CacheNameWithWebUser = "WebUser_"   //网站登录用户信息缓存
	CacheTimeWithWebUser = 86400 * 1000 //网站登录用户信息缓存时间

	CacheNameWithSystemUser = "SystemUser_" //系统登录用户信息缓存
	CacheTimeWithSystemUser = 86400 * 1000  //系统登录用户信息缓存时间

	CacheNameWebSocketUser = "WebSocketUser_" //websocket连接用户缓存
	CacheTimeWebSocketUser = 86400 * 1000 * 7 //websocket连接用户缓存时间
)
