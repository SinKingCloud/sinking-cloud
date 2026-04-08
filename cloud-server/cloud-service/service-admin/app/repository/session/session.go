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
	DeleteByUserIdAndTokens(userId int64, tokens []string) error
	DeleteByUserId(userId int64, tx ...*gorm.DB) error
	FindByToken(token string) (*model.Session, error)
	Select(where *SelectSession, orderByField string, orderByType string, page int, pageSize int) (list []*model.Session, total int64, err error)
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
