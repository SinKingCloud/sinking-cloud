package service

// ConfigService service接口
type ConfigService interface {
	GetUserConfigs(userId int64, group string, refresh bool) map[string]string //获取用户配置
	SetUserConfig(userId int64, group string, key string, value string) bool   //修改用户配置
	SetUserConfigs(userId int64, group string, data map[string]string)         //批量修改用户配置
}
