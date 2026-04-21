package scene

import (
	"encoding/json"
	"errors"
	"service-driver/app/constant"
	"service-driver/app/enum/scene_status"
	"service-driver/app/model"
	"service-driver/app/repository/scene"
	"strconv"
	"strings"
)

// Select 查询数据
func (s *service) Select(where *scene.SelectScene, orderByField string, orderByType string, page int, pageSize int) (list []*model.Scene, total int64, err error) {
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("状态参数错误")
		}
		if _, ok := scene_status.Map()[status]; !ok {
			return nil, 0, errors.New("状态参数不合法")
		}
	}
	if where.ChannelId != "" {
		_, err = strconv.Atoi(where.ChannelId)
		if err != nil {
			return nil, 0, errors.New("通道参数错误")
		}
	}
	return s.repositoryScene.Select(where, orderByField, orderByType, page, pageSize)
}

// GetIdNameMap 获取id name映射
func (s *service) GetIdNameMap(refresh bool) (map[int64]string, error) {
	cacheKey := constant.CacheNameWithSceneNameEnum
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
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithSceneNameEnum)
	}
	return maps, nil
}

// DeleteById 通过ID删除数据
func (s *service) DeleteById(id int64) error {
	exists, err := s.repositoryFile.ExistsBySceneId(id)
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
	s.cache.Delete(constant.CacheNameWithSceneNameEnum, constant.CacheNameWithScene+strconv.FormatInt(id, 10))
	return nil
}

// UpdateById 通过ID更新数据
func (s *service) UpdateById(id int64, data *scene.UpdateScene) error {
	if id <= 0 {
		return errors.New("ID不合法")
	}
	update := &scene.UpdateScene{}
	if data.Name != nil {
		if _, ok := data.Name.(string); !ok {
			return errors.New("场景名称不合法")
		}
		update.Name = data.Name
	}
	if data.AllowExt != nil {
		if _, ok := data.AllowExt.(string); !ok {
			return errors.New("允许类型不合法")
		}
		update.AllowExt = data.AllowExt
	}
	if data.Remark != nil {
		if _, ok := data.Remark.(string); !ok {
			return errors.New("场景备注不合法")
		}
		update.Remark = data.Remark
	}
	if data.MaxSize != nil {
		if _, ok := data.MaxSize.(int64); !ok {
			return errors.New("文件最大限制不合法")
		}
		update.MaxSize = data.MaxSize
	}
	if data.Status != nil {
		if value, ok := data.Status.(int); ok {
			if _, ok = scene_status.Map()[value]; !ok {
				return errors.New("场景状态不合法")
			}
		} else {
			return errors.New("场景状态不合法")
		}
		update.Status = data.Status
	}
	err := s.repositoryScene.UpdateById(id, data)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithSceneNameEnum, constant.CacheNameWithScene+strconv.FormatInt(id, 10))
	return nil
}

// Create 创建数据
func (s *service) Create(data *model.Scene) error {
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("场景名称不能为空")
	}
	if _, ok := scene_status.Map()[data.Status]; !ok {
		return errors.New("场景状态字段不合法")
	}
	if data.ChannelId <= 0 {
		return errors.New("存储通道字段不合法")
	}
	if data.MaxSize <= 0 {
		data.MaxSize = 0
	}
	if strings.TrimSpace(data.AllowExt) == "" {
		data.AllowExt = "*"
	}
	if info, err := s.repositoryChannel.FindById(data.Id); err != nil || info == nil {
		return errors.New("存储通道不存在")
	}
	err := s.repositoryScene.Create(data)
	if err != nil {
		return errors.Join(errors.New("创建失败"), err)
	}
	s.cache.Delete(constant.CacheNameWithSceneNameEnum)
	return nil
}

// FindById 通过ID查询数据
func (s *service) FindById(id int64, refresh bool) (*model.Scene, error) {
	cacheKey := constant.CacheNameWithScene + strconv.FormatInt(id, 10)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists && data != "" {
			var result *model.Scene
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	info, err := s.repositoryScene.FindById(id)
	if err != nil {
		return nil, err
	}
	if bytes, err2 := json.Marshal(info); err2 == nil {
		s.cache.SetWithExpire(cacheKey, string(bytes), constant.CacheTimeWithScene)
	}
	return info, nil
}
