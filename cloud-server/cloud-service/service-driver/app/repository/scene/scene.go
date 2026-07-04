package scene

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Scene) error                                                            //创建数据
	DeleteById(id int64) error                                                                 //根据ID删除数据
	UpdateById(id int64, data *UpdateScene) error                                              //根据id更新数据
	FindById(id int64) (*model.Scene, error)                                                   //根据ID查询数据
	ExistsByChannelId(channelId int64) (bool, error)                                           //判断channelId是否被使用
	SelectIdNameMap() (map[int64]string, error)                                                //查询ID name
	Select(where *SelectScene, page *page.Query) (list []*model.Scene, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository(db, &model.Scene{})}
}
