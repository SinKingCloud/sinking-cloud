package file

import (
	"service-driver/app/model"
	"service-driver/app/repository/file"
)

// Select 查询数据
func (s *service) Select(where *file.SelectFile, orderByField string, orderByType string, page int, pageSize int) (list []*model.File, total int64, err error) {
	return s.repositoryFile.Select(where, orderByField, orderByType, page, pageSize)
}
