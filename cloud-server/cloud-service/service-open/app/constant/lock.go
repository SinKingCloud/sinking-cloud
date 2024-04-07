package constant

const (
	LockOauthAuthTicket     = "LockOauthAuthTicket_" //授权登录并发锁
	LockTimeOauthAuthTicket = 1 * 60 * 1000          //授权登录并发锁最大时间

	LockAuthBuy     = "LockAuthBuy_" //授权应用下单并发锁
	LockTimeAuthBuy = 1 * 60 * 1000  //授权应用下单并发锁最大时间
)
