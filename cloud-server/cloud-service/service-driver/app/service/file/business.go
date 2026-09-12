package file

import (
	"context"
	"encoding/json"
	"errors"
	"service-driver/app/constant"
	"service-driver/app/enum/file_status"
	"service-driver/app/model"
	"service-driver/app/repository/file"
	"strconv"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Select 查询数据
func (s *service) Select(where *file.SelectFile, query *page.Query) (*page.Result[*model.File], error) {
	if where.Status != nil {
		if _, ok := file_status.Map()[*where.Status]; !ok {
			return nil, errors.New("状态参数不合法")
		}
	}
	return s.repositoryFile.Select(where, query)
}

// GetUrlByHash 根据 hash 获取文件公开地址。
func (s *service) GetUrlByHash(sceneId int64, hash string) (string, error) {
	if sceneId <= 0 {
		return "", errors.New("文件场景不能为空")
	}
	if hash == "" {
		return "", errors.New("文件哈希不能为空")
	}
	const cacheErrorPrefix = "error:"
	cacheKey := s.buildPreviewCacheKey(sceneId, hash)
	if data, exists := s.cache.Get(cacheKey); exists {
		if data == "" {
			return "", errors.New("文件不存在")
		}
		if strings.HasPrefix(data, cacheErrorPrefix) {
			return "", errors.New(strings.TrimPrefix(data, cacheErrorPrefix))
		}
		return data, nil
	}
	cacheError := func(err error) (string, error) {
		s.cache.SetWithExpire(cacheKey, cacheErrorPrefix+err.Error(), constant.CacheTimeWithFilePreview)
		return "", err
	}
	info, err := s.findByHashWithSceneId(hash, sceneId, false)
	if err != nil {
		return cacheError(errors.Join(errors.New("获取文件信息失败"), err))
	}
	if info == nil {
		return cacheError(errors.New("文件不存在"))
	}
	if info.Status != file_status.Enable {
		return cacheError(errors.New("文件不可访问"))
	}
	if info.ChannelId <= 0 {
		return cacheError(errors.New("文件通道不存在"))
	}
	if info.Path == "" {
		return cacheError(errors.New("文件路径不存在"))
	}
	client, err := s.channelService.GetStorage(info.ChannelId)
	if err != nil {
		return cacheError(errors.Join(errors.New("获取存储信息失败"), err))
	}
	url, err := client.PublicURL(info.Path)
	if err != nil {
		return cacheError(errors.Join(errors.New("获取文件访问地址失败"), err))
	}
	s.cache.SetWithExpire(cacheKey, url, constant.CacheTimeWithFilePreview)
	return url, nil
}

// UpdateByIds 批量更新数据
func (s *service) UpdateByIds(ids []int64, data *file.UpdateFile) error {
	if len(ids) == 0 || data == nil {
		return errors.New("参数错误")
	}
	list, err := s.repositoryFile.SelectChannelIdSceneIdHashPathByIds(ids)
	if err != nil {
		return errors.Join(errors.New("获取文件信息失败"), err)
	}
	update := &file.UpdateFile{}
	if data.Status != nil {
		if _, ok := file_status.Map()[*data.Status]; !ok {
			return errors.New("文件状态不合法")
		}
		update.Status = data.Status
	}
	err = s.repositoryFile.UpdateByIds(ids, update)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	s.clearCache(list)
	return nil
}

// DeleteByIds 批量删除数据
func (s *service) DeleteByIds(ids []int64, clear bool) error {
	if len(ids) == 0 {
		return errors.New("参数错误")
	}
	list, err := s.repositoryFile.SelectChannelIdSceneIdHashPathByIds(ids)
	if err != nil {
		return errors.Join(errors.New("获取文件信息失败"), err)
	}
	var paths map[int64][]string
	if clear {
		paths = make(map[int64][]string)
		for _, item := range list {
			if item == nil || item.ChannelId <= 0 || item.Path == "" {
				continue
			}
			paths[item.ChannelId] = append(paths[item.ChannelId], item.Path)
		}
	}
	err = s.repositoryFile.DeleteByIds(ids)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	s.clearCache(list)
	for k, v := range paths {
		if client, err2 := s.channelService.GetStorage(k); err2 == nil && client != nil {
			_ = client.DeleteMany(context.Background(), v)
		}
	}
	return nil
}

// findByHashWithSceneId 根据 hash 和场景查询文件，并优先使用缓存。
func (s *service) findByHashWithSceneId(hash string, sceneId int64, refresh bool) (*model.File, error) {
	if sceneId <= 0 || hash == "" {
		return nil, nil
	}
	cacheKey := s.buildFileCacheKey(sceneId, hash)
	if !refresh {
		if data, exists := s.cache.Get(cacheKey); exists {
			if data == "" {
				return nil, nil
			}
			var result *model.File
			if err := json.Unmarshal([]byte(data), &result); err == nil {
				return result, nil
			}
		}
	}
	info, err := s.repositoryFile.FindByHashWithSceneId(hash, sceneId)
	if err != nil {
		return nil, err
	}
	if data, err := json.Marshal(info); err == nil {
		s.cache.SetWithExpire(cacheKey, string(data), constant.CacheTimeWithFile)
	}
	return info, nil
}

// clearCache 清理文件缓存。
func (s *service) clearCache(list []*model.File) {
	if len(list) == 0 {
		return
	}
	keys := make([]string, 0, len(list))
	seen := make(map[string]struct{}, len(list))
	for _, item := range list {
		if item == nil || item.SceneId <= 0 || item.Hash == "" {
			continue
		}
		cacheKeys := []string{
			s.buildFileCacheKey(item.SceneId, item.Hash),
			s.buildPreviewCacheKey(item.SceneId, item.Hash),
		}
		for _, cacheKey := range cacheKeys {
			if _, exists := seen[cacheKey]; exists {
				continue
			}
			seen[cacheKey] = struct{}{}
			keys = append(keys, cacheKey)
		}
	}
	if len(keys) > 0 {
		s.cache.Delete(keys...)
	}
}

// buildFileCacheKey 生成文件缓存key。
func (s *service) buildFileCacheKey(sceneId int64, hash string) string {
	return constant.CacheNameWithFile + strconv.FormatInt(sceneId, 10) + "_" + hash
}

// buildPreviewCacheKey 生成文件预览缓存key。
func (s *service) buildPreviewCacheKey(sceneId int64, hash string) string {
	return constant.CacheNameWithFilePreview + strconv.FormatInt(sceneId, 10) + "_" + hash
}
