package repository

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"

// Repository 基础仓储，T 表示分页查询返回的数据类型。
type Repository[T any] struct {
	Database *database.Database
}

// NewRepository 创建基础仓储
func NewRepository[T any](db *database.Database) *Repository[T] {
	return &Repository[T]{Database: db}
}
