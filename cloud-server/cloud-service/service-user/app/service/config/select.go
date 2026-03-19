package config

// selectWithUserIDAndKey 获取配置数据
func (s *service) selectWithUserIDAndKey(userId int64, key string) map[string]string {
	configs, err := s.repositoryConfig.FindByUserIdAndKey(userId, key)
	temp := make(map[string]string)
	if err != nil {
		return temp
	}
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}
