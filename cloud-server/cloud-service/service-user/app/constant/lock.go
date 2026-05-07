package constant

const (
	LockConfigSet     = "ConfigSet_"  //修改配置并发锁
	LockTimeConfigSet = 1 * 60 * 1000 //修改配置并发锁最大时间

	LockUserCreate     = "UserCreate_" //修改配置并发锁
	LockTimeUserCreate = 1 * 60 * 1000 //修改配置并发锁最大时间

	LockRealAuth     = "RealAuth_"   //实名认证并发锁
	LockTimeRealAuth = 5 * 60 * 1000 //实名认证并发锁最大时间
)
