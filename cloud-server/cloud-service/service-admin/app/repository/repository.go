package repository

import (
	"service-admin/app/model"
	"service-admin/app/repository/config"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Config 配置仓储接口
type Config interface {
	FindByIdentIdAndKey(identId, key string) ([]*model.Config, error)
	UpdateByIdentIdAndKey(identId, key, value string) (int64, error)
	Create(config *model.Config) error
}

// NewConfig 创建Config仓储
func NewConfig(db *database.Database) *config.Repository {
	return &config.Repository{Database: db}
}
