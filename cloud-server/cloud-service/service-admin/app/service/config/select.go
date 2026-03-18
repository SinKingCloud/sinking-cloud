package config

// selectWithIdentIDAndKey 获取配置数据
func (s *service) selectWithIdentIDAndKey(identId string, key string) map[string]string {
	configs, err := s.repo.FindByIdentIdAndKey(identId, key)
	temp := make(map[string]string)
	if err != nil {
		return temp
	}
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}
