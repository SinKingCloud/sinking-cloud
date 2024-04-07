package constant

const (
	LockUinServer     = "LockUinServer_" //QQ添加锁
	LockTimeUinServer = 1 * 60 * 1000    //QQ添加锁

	LockUinJob     = "LockUinJob_"  //挂机任务锁
	LockTimeUinJob = 10 * 60 * 1000 //挂机任务锁最大时间

	LockUinCheckJob     = "LockUinCheckJob_" //挂机检测任务锁
	LockTimeUinCheckJob = 10 * 60 * 1000     //挂机检测任务锁最大时间

	LockUserUinJob     = "LockUserUinJob_" //用户挂机任务锁
	LockTimeUserUinJob = 1 * 60 * 1000     //用户挂机任务锁等待时间

	LockUserUinOperate     = "LockUserUinOperate_" //用户挂机操作锁
	LockTimeUserUinOperate = 5 * 1000              //用户挂机操作锁等待时间

	LockUinBuy     = "LockUinBuy_" //挂机应用下单并发锁
	LockTimeUinBuy = 1 * 60 * 1000 //挂机应用下单并发锁最大时间

	LockUinCookie     = "LockUinCookie_" //挂机应用下单并发锁
	LockTimeUinCookie = 1 * 60 * 1000    //挂机应用下单并发锁最大时间

	LockUinInfo     = "LockUinInfo_" //QQ更新详情并发锁
	LockTimeUinInfo = 1 * 60 * 1000  //QQ更新详情并发锁最大时间
)
