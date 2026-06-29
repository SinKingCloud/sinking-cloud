package message

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Interface 消息仓储接口
type Interface interface {
	Create(message *model.Message) error
	Creates(messages []*model.Message, tx ...*gorm.DB) error
	DeleteByIdsAndUserId(ids []int64, userId int64, tx ...*gorm.DB) error
	DeleteByUserId(userId int64, tx ...*gorm.DB) error
	FindByIdAndUserId(id int64, userId int64) (*model.Message, error)
	UpdateByIdAndUserId(id int64, userId int64, data *UpdateMessage, tx ...*gorm.DB) error
	ReadByIdsAndUserId(ids []int64, userId int64, readTime str.DateTime, tx ...*gorm.DB) error
	Select(where *SelectMessage, orderByField string, orderByType string, page int, pageSize int) (list []*Message, total int64, err error)
}

// Repository 消息仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.Message{})}
}
