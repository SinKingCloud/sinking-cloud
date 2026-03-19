package service

// ConfigService service接口
type ConfigService interface {
	GetUserConfigs(userId int64, group string, refresh bool) map[string]string
	SetUserConfig(userId int64, group string, key string, value string) bool
	SetUserConfigs(userId int64, group string, data map[string]string)
}
