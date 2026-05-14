package constant

const (
	LockConfigSet     = "LockConfigSet_" //修改配置并发锁
	LockTimeConfigSet = 1 * 60 * 1000    //修改配置并发锁最大时间

	LockPermission     = "LockPermission_" //权限注册并发锁
	LockTimePermission = 1 * 60 * 1000     //权限注册并发锁最大时间

	LockRole     = "LockRole_"   //角色业务并发锁
	LockTimeRole = 1 * 60 * 1000 //角色业务并发锁最大时间

	LockMenuOperate     = "LockMenuOperate" //菜单操作并发锁
	LockTimeMenuOperate = 1 * 60 * 1000     //菜单操作并发锁最大时间

	LockUserRoleAssign     = "LockUserRoleAssign_" //用户分配角色并发锁
	LockTimeUserRoleAssign = 1 * 60 * 1000         //用户分配角色并发锁最大时间
)
