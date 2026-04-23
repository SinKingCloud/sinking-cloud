package order_log

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 仓储接口
type Interface interface {
	FindById(id int64) (*model.OrderLog, error)                                                                                                     //根据ID查询记录
	Create(data *model.OrderLog) error                                                                                                              //创建数据
	Select(where *SelectOrderLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.OrderLog, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
