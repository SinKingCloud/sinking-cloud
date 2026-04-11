package constant

const (
	CacheNameWithSocketToken = "SocketToken_" //token授权
	CacheTimeWithSocketToken = 10 * 60 * 1000 //token授权储存时间

	CacheNameWithSocketAddr            = "SocketAddr_"  //token绑定地址
	CacheTimeWithSocketAddr            = 60 * 60 * 1000 //token绑定地址储存时间
	CacheTimeWithSocketAddrRenewMinute = 5              //token绑定地址续期间隔(分钟)
)
