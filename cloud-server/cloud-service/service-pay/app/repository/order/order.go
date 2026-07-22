package order

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Order, tx ...*gorm.DB) error                             //创建数据
	FindById(id int64) (*model.Order, error)                                    //根据ID查询记录
	CheckChannelHasUnPayOrder(channelId int64) (bool, error)                    //根据通道ID判断是否存在未过期且未支付订单
	FindExpiredUnPayIds(expiredTime string, limit int) ([]int64, error)         //查询过期未支付订单ID
	Select(where *SelectOrder, query *page.Query) (*page.Result[*Order], error) //查询数据
	UpdateById(id int64, data *UpdateOrder, tx ...*gorm.DB) error               //根据id更新数据
	DeleteByIds(ids []int64, tx ...*gorm.DB) error                              //根据ID批量删除数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*Order]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*Order](db)}
}
