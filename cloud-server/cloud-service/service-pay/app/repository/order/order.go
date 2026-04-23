package order

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Order) error                                                                                                     //创建数据
	FindByTradeNo(tradeNo string) (*model.Order, error)                                                                                 //根据订单号查询记录
	FindById(id int64) (*model.Order, error)                                                                                            //根据ID查询记录
	Select(where *SelectOrder, orderByField string, orderByType string, page int, pageSize int) (list []*Order, total int64, err error) //查询数据
	UpdateById(id int64, data *UpdateOrder) error                                                                                       //根据id更新数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
