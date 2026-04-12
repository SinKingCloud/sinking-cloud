package session

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface Session仓储接口
type Interface interface {
	Create(session *model.Session) error
	DeleteByToken(token string) error
	DeleteByUserIdAndTokens(userId int64, tokens []string) error
	FindByToken(token string) (*model.Session, error)
	Select(where *SelectSession, orderByField string, orderByType string, page int, pageSize int) (list []*model.Session, total int64, err error)
	SelectUnExpireTokensByUserIds(userIds []int64) ([]string, error)
}

// Repository Session仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
