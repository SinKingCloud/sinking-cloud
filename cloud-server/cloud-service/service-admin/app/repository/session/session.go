package session

import (
	"service-admin/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface Session仓储接口
type Interface interface {
	Create(session *model.Session) error
	DeleteByToken(token string) error
	DeleteByUserId(userId int64, tx ...*gorm.DB) error
	FindByToken(token string) (*model.Session, error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
