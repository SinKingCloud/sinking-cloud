package repository

import (
	"service-user/app/model"
	"service-user/app/repository/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Config 配置仓储接口
type Config interface {
	FindByUserIdAndKey(userId int64, key string) ([]*model.Config, error)
	UpdateByUserIdAndKey(userId int64, key string, value string) (int64, error)
	Create(config *model.Config) error
}

// NewConfig 创建Config仓储
func NewConfig(db *database.Database) *config.Repository {
	return &config.Repository{Database: db}
}
