package config

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 配置仓储接口
type Interface interface {
	FindByIdentIdAndKey(identId, key string) ([]*model.Config, error)
	UpdateByIdentIdAndKey(identId, key, value string) (int64, error)
	DeleteByIdentId(identId string, tx ...*gorm.DB) error
	Create(config *model.Config) error
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
