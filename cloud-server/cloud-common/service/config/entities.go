package config

// GetWebUserConfigParam 获取Web用户配置请求参数
type GetWebUserConfigParam struct {
	UserId  int64  // 用户ID
	Group   string // 配置分组
	Refresh bool   // 是否强制刷新缓存
}

// SetWebUserConfigParam 设置Web用户配置请求参数
type SetWebUserConfigParam struct {
	UserId int64  // 用户ID
	Key    string // 配置键
	Value  string // 配置值
}

// SetWebUserConfigsParam 批量设置Web用户配置请求参数
type SetWebUserConfigsParam struct {
	UserId int64             // 用户ID
	Data   map[string]string // 配置数据
}

// GetSystemUserConfigParam 获取系统用户配置请求参数
type GetSystemUserConfigParam struct {
	UserId  int64  // 用户ID
	Group   string // 配置分组
	Refresh bool   // 是否强制刷新缓存
}

// SetSystemUserConfigParam 设置系统用户配置请求参数
type SetSystemUserConfigParam struct {
	UserId int64  // 用户ID
	Key    string // 配置键
	Value  string // 配置值
}

// SetSystemUserConfigsParam 批量设置系统用户配置请求参数
type SetSystemUserConfigsParam struct {
	UserId int64             // 用户ID
	Data   map[string]string // 配置数据
}

// GetSystemConfigParam 获取系统配置请求参数
type GetSystemConfigParam struct {
	Group   string // 配置分组
	Refresh bool   // 是否强制刷新缓存
}

// GetSystemConfigsParam 批量获取系统配置请求参数
type GetSystemConfigsParam struct {
	Group   []string // 配置分组
	Refresh bool     // 是否强制刷新缓存
}

// SetSystemConfigParam 设置系统配置请求参数
type SetSystemConfigParam struct {
	Key   string // 配置键
	Value string // 配置值
}

// SetSystemConfigsParam 批量设置系统配置请求参数
type SetSystemConfigsParam struct {
	Data map[string]string // 配置数据
}
