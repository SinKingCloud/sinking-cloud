package order_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	FindById(id int64) (*model.OrderLog, error)                                                //根据ID查询记录
	Create(data *model.OrderLog, tx ...*gorm.DB) error                                         //创建数据
	DeleteByOrderIds(orderIds []int64, tx ...*gorm.DB) error                                   //根据订单ID批量删除数据
	Select(where *SelectOrderLog, page *page.Query) (list []*OrderLog, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.OrderLog{})}
}
