package file

import (
	"service-driver/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/repository"
)

// Interface 仓储接口
type Interface interface {
	Create(data *model.File) error                                                   //创建数据
	DeleteByIds(ids []int64) error                                                   //根据ids删除数据
	ExistsByChannelId(channelId int64) (bool, error)                                 //判断channelId是否被使用
	ExistsBySceneId(sceneId int64) (bool, error)                                     //判断sceneId是否被使用
	FindByHashWithSceneId(hash string, sceneId int64) (*model.File, error)           //根据hash和场景查询文件
	SelectChannelIdSceneIdHashPathByIds(ids []int64) (list []*model.File, err error) //根据ids查询清理所需的通道、场景、哈希和路径
	Select(where *SelectFile, query *page.Query) (*page.Result[*model.File], error)  //查询数据
	UpdateByIds(ids []int64, perm *UpdateFile) error                                 //根据 ids 更新数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.File]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.File](db)}
}
