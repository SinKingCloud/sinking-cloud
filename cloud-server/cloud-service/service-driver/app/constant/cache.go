package constant

const (
	CacheNameWithChannel = "Channel_"        //通道缓存
	CacheTimeWithChannel = int64(600 * 1000) //通道缓存储存时间

	CacheNameWithScene = "Scene_"          //场景缓存
	CacheTimeWithScene = int64(600 * 1000) //场景缓存储存时间

	CacheNameWithFile = "File_"                        //文件缓存
	CacheTimeWithFile = int64(7 * 24 * 60 * 60 * 1000) // 文件缓存储存时间

	CacheNameWithFilePreview = "FilePreview_"                 //文件预览缓存
	CacheTimeWithFilePreview = int64(7 * 24 * 60 * 60 * 1000) // 文件预览缓存储存时间

	CacheNameWithChannelNameEnum = "ChannelNameEnum" //枚举缓存
	CacheTimeWithChannelNameEnum = int64(600 * 1000) //枚举缓存储存时间

	CacheNameWithSceneNameEnum = "SceneNameEnum"   //枚举缓存
	CacheTimeWithSceneNameEnum = int64(600 * 1000) //枚举缓存储存时间
)
