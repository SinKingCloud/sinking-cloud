package channel

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Channel) error                                               //创建数据
	DeleteById(id int64) error                                                      //根据ID删除数据
	SelectIdNameMap() (map[int64]string, error)                                     //查询ID name
	UpdateById(id int64, perm *UpdateChannel) error                                 //根据id更新
	FindById(id int64) (*model.Channel, error)                                      //根据ID查询
	Select(where *SelectChannel, query *page.Query) (*page.Result[*Channel], error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*Channel]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*Channel](db)}
}
