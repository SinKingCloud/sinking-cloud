package file

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.File) error                                                                                                          //创建数据
	DeleteByIds(ids []int64) error                                                                                                          //根据ids删除数据
	ExistsByChannelId(channelId int64) (bool, error)                                                                                        //判断channelId是否被使用
	ExistsBySceneId(sceneId int64) (bool, error)                                                                                            //判断sceneId是否被使用
	Select(where *SelectFile, orderByField string, orderByType string, page int, pageSize int) (list []*model.File, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
