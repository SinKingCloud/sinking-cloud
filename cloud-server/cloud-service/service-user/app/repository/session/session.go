package session

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface Session仓储接口
type Interface interface {
	Create(session *model.Session) error
	DeleteByToken(token string) error
	FindByToken(token string) (*model.Session, error)
}

// Repository Session仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
