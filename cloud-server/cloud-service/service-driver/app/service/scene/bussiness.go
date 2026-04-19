package scene

import (
	"encoding/json"
	"errors"
	"service-driver/app/constant"
	"service-driver/app/model"
	"service-driver/app/repository/scene"
)

// Select 查询数据
func (s *service) Select(where *scene.SelectScene, orderByField string, orderByType string, page int, pageSize int) (list []*model.Scene, total int64, err error) {
	return s.repositoryScene.Select(where, orderByField, orderByType, page, pageSize)
}

// GetIdNameMap 获取id name映射
func (s *service) GetIdNameMap(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithEnum + "scene_names"
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result map[int64]string
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	maps, err := s.repositoryScene.SelectIdNameMap()
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
	exists, err := s.repositoryFile.ExistsByChannelId(id)
	if err != nil {
		return errors.New("获取文件使用信息失败")
	}
	if exists {
		return errors.New("该场景下的文件数据仍在使用,请先清理文件数据")
	}
	err = s.repositoryScene.DeleteById(id)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	return nil
}
