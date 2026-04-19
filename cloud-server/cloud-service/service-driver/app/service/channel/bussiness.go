package channel

import (
	"encoding/json"
	"errors"
	"service-driver/app/constant"
	"service-driver/app/repository/channel"
)

// Select 查询数据
func (s *service) Select(where *channel.SelectChannel, orderByField string, orderByType string, page int, pageSize int) (list []*channel.Channel, total int64, err error) {
	return s.repositoryChannel.Select(where, orderByField, orderByType, page, pageSize)
}

// GetIdNameMap 获取id name映射
func (s *service) GetIdNameMap(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithEnum + "channel_names"
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result map[int64]string
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	maps, err := s.repositoryChannel.SelectIdNameMap()
	if err != nil {
		return nil, err
	}
	if data, err2 := json.Marshal(maps); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithEnum)
	}
	return maps, nil
}

// DeleteById 通过ID删除数据
func (s *service) DeleteById(id int64) error {
	exists, err := s.repositoryScene.ExistsByChannelId(id)
	if err != nil {
		return errors.New("获取场景使用信息失败")
	}
	if exists {
		return errors.New("该通道下的场景数据仍在使用,请先清理场景数据")
	}
	exists, err = s.repositoryFile.ExistsByChannelId(id)
	if err != nil {
		return errors.New("获取文件使用信息失败")
	}
	if exists {
		return errors.New("该通道下的文件数据仍在使用,请先清理文件数据")
	}
	err = s.repositoryChannel.DeleteById(id)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	return nil
}
