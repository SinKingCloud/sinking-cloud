package balance_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.BalanceLog, tx ...*gorm.DB) error                                                 //创建数据
	FindById(id int64) (*model.BalanceLog, error)                                                        //通过ID查询数据
	Select(where *SelectBalanceLog, page *page.Query) (list []*model.BalanceLog, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.BalanceLog{})}
}
