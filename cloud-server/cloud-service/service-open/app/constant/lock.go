package constant

const (
	LockUinServer     = "LockUinServer_" //QQ添加锁
	LockTimeUinServer = 1 * 60 * 1000    //QQ添加锁

	LockUinServerChange     = "LockUinServerChange_" //挂机服务器故障迁移锁
	LockTimeUinServerChange = 10 * 60 * 1000         //挂机服务器故障迁移锁最大时间

	LockUinServerCheck     = "LockUinServerCheck_" //挂机服务器socks5代理检测锁
	LockTimeUinServerCheck = 10 * 60 * 1000        //挂机服务器socks5代理检测任务锁最大时间

	LockUinJob     = "LockUinJob_"  //挂机任务锁
	LockTimeUinJob = 10 * 60 * 1000 //挂机任务锁最大时间

	LockUinCheckJob     = "LockUinCheckJob_" //挂机检测任务锁
	LockTimeUinCheckJob = 10 * 60 * 1000     //挂机检测任务锁最大时间

	LockUserUinJob     = "LockUserUinJob_" //用户挂机任务锁
	LockTimeUserUinJob = 5 * 60 * 1000     //用户挂机任务锁等待时间

	LockUserUinOperate     = "LockUserUinOperate_" //用户挂机操作锁
	LockTimeUserUinOperate = 5 * 1000              //用户挂机操作锁等待时间

	LockUinBuy     = "LockUinBuy_" //挂机应用下单并发锁
	LockTimeUinBuy = 1 * 60 * 1000 //挂机应用下单并发锁最大时间

	LockUinCookie     = "LockUinCookie_" //挂机应用下单并发锁
	LockTimeUinCookie = 1 * 60 * 1000    //挂机应用下单并发锁最大时间

	LockOauthAuthTicket     = "LockOauthAuthTicket_" //授权登录并发锁
	LockTimeOauthAuthTicket = 1 * 60 * 1000          //授权登录并发锁最大时间

	LockUinInfo     = "LockUinInfo_" //QQ更新详情并发锁
	LockTimeUinInfo = 1 * 60 * 1000  //QQ更新详情并发锁最大时间

	LockAuthBuy     = "LockAuthBuy_" //授权应用下单并发锁
	LockTimeAuthBuy = 1 * 60 * 1000  //授权应用下单并发锁最大时间
)
