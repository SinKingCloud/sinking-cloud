package constant

const (
	LockTask     = "LockTask_"    // 任务执行并发锁
	LockTimeTask = 60 * 60 * 1000 // 任务执行并发锁最大时间

	LockTaskDispatch     = "LockTaskDispatch_" // 任务派发并发锁
	LockTimeTaskDispatch = 10 * 1000           // 任务派发并发锁最大时间

	LockTaskDispatchRenew     = "LockTaskDispatchRenew_" // 任务投递节点续期锁
	LockTimeTaskDispatchRenew = 10 * 1000                // 任务投递节点续期锁最大时间
)
