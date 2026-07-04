package session

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface Session仓储接口
type Interface interface {
	Create(session *model.Session) error
	DeleteByToken(token string) error
	DeleteByUserIdAndTokens(userId int64, tokens []string) error
	DeleteByIds(ids []int64, tx ...*gorm.DB) error
	FindByToken(token string) (*model.Session, error)
	SelectExpireSessions(tx ...*gorm.DB) (map[int64]string, error)
	Select(where *SelectSession, page *page.Query) (list []*model.Session, total int64, err error)
	SelectUnExpireTokensByUserIds(userIds []int64) ([]string, error)
}

// Repository Session仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.Session{})}
}
