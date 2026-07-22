package constant

const (
	LockMindDataOperate     = "MindDataOperate_" //导图操作并发锁
	LockTimeMindDataOperate = 1 * 60 * 1000      //导图操作并发锁最大时间
	LockMindDataCreate      = "MindDataCreate_"  //创建导图并发锁
	LockTimeMindDataCreate  = 1 * 60 * 1000      //创建导图并发锁最大时间

	LockMindTypeOperate     = "MindTypeOperate_" //导图分类操作并发锁
	LockTimeMindTypeOperate = 1 * 60 * 1000      //导图分类操作并发锁最大时间
)
