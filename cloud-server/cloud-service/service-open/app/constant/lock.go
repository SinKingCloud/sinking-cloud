package constant

const (
	LockProxySocks5Check     = "LockUinJob_"  //socks5代理检测锁
	LockTimeProxySocks5Check = 10 * 60 * 1000 //socks5代理检测任务锁最大时间

	LockUinJob     = "LockUinJob_"  //挂机任务锁
	LockTimeUinJob = 10 * 60 * 1000 //挂机任务锁最大时间

	LockUinCheckJob     = "LockUinCheckJob_" //挂机检测任务锁
	LockTimeUinCheckJob = 10 * 60 * 1000     //挂机检测任务锁最大时间

	LockUserUinJob     = "LockUserUinJob_" //用户挂机任务锁
	LockTimeUserUinJob = 5 * 60 * 1000     //用户挂机任务锁等待时间
)
