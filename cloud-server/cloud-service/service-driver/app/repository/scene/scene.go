package scene

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.Scene) error                                                                                                           //创建数据
	DeleteById(id int64, tx ...*gorm.DB) error                                                                                                //根据ID删除数据
	UpdateByIds(ids []int64, perm *UpdateScene) error                                                                                         //根据ids更新数据
	FindById(id int64) (*model.Scene, error)                                                                                                  //根据ID查询数据
	SelectIdNameMap() (map[int64]string, error)                                                                                               //查询ID name
	Select(where *SelectScene, orderByField string, orderByType string, page int, pageSize int) (list []*model.Scene, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
