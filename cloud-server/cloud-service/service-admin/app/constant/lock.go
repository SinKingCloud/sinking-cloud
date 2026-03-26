package constant

const (
	LockConfigSet     = "LockConfigSet_" //修改配置并发锁
	LockTimeConfigSet = 1 * 60 * 1000    //修改配置并发锁最大时间

	LockPermission     = "LockPermission_" //权限注册并发锁
	LockTimePermission = 1 * 60 * 1000     //权限注册并发锁最大时间
)
