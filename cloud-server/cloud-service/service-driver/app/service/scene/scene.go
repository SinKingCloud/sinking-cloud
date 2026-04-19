package scene

import (
	"service-driver/app/model"
	"service-driver/app/repository/file"
	"service-driver/app/repository/scene"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Service service接口
type Service interface {
	Select(where *scene.SelectScene, orderByField string, orderByType string, page int, pageSize int) (list []*model.Scene, total int64, err error) //查询数据
	GetIdNameMap(refresh bool) (map[int64]string, error)                                                                                            //获取id name映射
	DeleteById(id int64) error                                                                                                                      //通过ID删除数据
}

// Session服务
type service struct {
	repositoryScene scene.Interface
	repositoryFile  file.Interface
	database        *database.Database
	cache           cache.Interface
}

// NewService 实例化service
func NewService(repositoryScene scene.Interface, repositoryFile file.Interface, database *database.Database, cache cache.Interface) *service {
	return &service{
		repositoryScene: repositoryScene,
		repositoryFile:  repositoryFile,
		database:        database,
		cache:           cache,
	}
}
