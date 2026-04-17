package constant

const (
	CacheNameWithTaskInfo = "TaskInfo_"       //任务信息
	CacheTimeWithTaskInfo = int64(600 * 1000) //任务信息缓存时间

	CacheNameWithTaskVersion = "TaskVersion"                   //任务版本
	CacheTimeWithTaskVersion = int64(30 * 24 * 60 * 60 * 1000) //任务版本缓存时间
)
