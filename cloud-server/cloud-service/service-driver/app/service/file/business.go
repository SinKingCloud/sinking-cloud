package file

import (
	"context"
	"errors"
	"service-driver/app/enum/file_status"
	"service-driver/app/model"
	"service-driver/app/repository/file"
	"strconv"
)

// Select 查询数据
func (s *service) Select(where *file.SelectFile, orderByField string, orderByType string, page int, pageSize int) (list []*model.File, total int64, err error) {
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("状态参数错误")
		}
		if _, ok := file_status.Map()[status]; !ok {
			return nil, 0, errors.New("状态参数不合法")
		}
	}
	return s.repositoryFile.Select(where, orderByField, orderByType, page, pageSize)
}

// UpdateByIds 批量更新数据
func (s *service) UpdateByIds(ids []int64, data *file.UpdateFile) error {
	if len(ids) == 0 || data == nil {
		return errors.New("参数错误")
	}
	update := &file.UpdateFile{}
	if data.Status != nil {
		if value, ok := data.Status.(int); ok {
			if _, ok = file_status.Map()[value]; !ok {
				return errors.New("文件状态不合法")
			}
		} else {
			return errors.New("文件状态不合法")
		}
		update.Status = data.Status
	}
	err := s.repositoryFile.UpdateByIds(ids, update)
	if err != nil {
		return errors.Join(errors.New("更新失败"), err)
	}
	return nil
}

// DeleteByIds 批量删除数据
func (s *service) DeleteByIds(ids []int64, clear bool) error {
	if len(ids) == 0 {
		return errors.New("参数错误")
	}
	var paths map[int64][]string
	if clear {
		list, err := s.repositoryFile.SelectChannelIdWithPathByIds(ids)
		if err != nil {
			return errors.Join(errors.New("获取清理路径失败"), err)
		}
		paths = make(map[int64][]string)
		for _, item := range list {
			if item == nil || item.ChannelId <= 0 || item.Path == "" {
				continue
			}
			paths[item.ChannelId] = append(paths[item.ChannelId], item.Path)
		}
	}
	err := s.repositoryFile.DeleteByIds(ids)
	if err != nil {
		return errors.Join(errors.New("删除失败"), err)
	}
	for k, v := range paths {
		if client, err2 := s.channelService.GetStorage(k); err2 == nil && client != nil {
			_ = client.DeleteMany(context.Background(), v)
		}
	}
	return nil
}
