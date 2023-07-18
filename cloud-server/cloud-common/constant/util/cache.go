package util

const (
	CacheNameWithSystemConfig = "SystemConfig" //系统配置缓存
	CacheTimeWithSystemConfig = 86400 * 1000   //默认系统配置缓存时间

	CacheNameWithWebUser = "WebUser_"   //网站登录用户信息缓存
	CacheTimeWithWebUser = 86400 * 1000 //默认网站登录用户信息缓存时间

	CacheNameWithSystemUser = "SystemUser_" //系统登录用户信息缓存
	CacheTimeWithSystemUser = 86400 * 1000  //默认系统登录用户信息缓存时间

	CacheNameWebSocketAuth = "WebSocketAuth_" //websocket token缓存
	CacheTimeWebSocketAuth = 600 * 1000       //默认websocket token缓存时间

	CacheNameWebSocketToken = "WebSocketToken_" //websocket连接用户缓存
	CacheTimeWebSocketToken = 86400 * 1000 * 7  //默认websocket连接用户缓存时间
)
