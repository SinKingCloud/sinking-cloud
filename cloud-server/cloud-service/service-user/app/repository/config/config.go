package config

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
)

// Interface 配置仓储接口
type Interface interface {
	FindByUserIdAndKey(userId int64, key string) ([]*model.Config, error)
	UpdateByUserIdAndKey(userId int64, key string, value string) (int64, error)
	Create(config *model.Config) error
}

// Repository 配置仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.Config{})}
}
