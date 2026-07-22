package channel

import (
	"service-driver/app/repository/channel"
	"service-driver/app/repository/file"
	"service-driver/app/repository/scene"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/storage"
)

// Service service接口
type Service interface {
	Create(name string, types int, status int, config string) error                                 //创建数据
	Select(where *channel.SelectChannel, query *page.Query) (*page.Result[*channel.Channel], error) //查询数据
	GetIdNameMap(refresh bool) (map[int64]string, error)                                            //获取id name映射
	DeleteById(id int64) error                                                                      //通过ID删除数据
	UpdateById(id int64, data *channel.UpdateChannel) error                                         //通过ID更新数据
	FindById(id int64, refresh bool) (*Channel, error)                                              //通过ID查询数据
	GetStorage(id int64) (storage.Storage, error)                                                   //获取储存对象
}

// Session服务
type service struct {
	repositoryChannel channel.Interface
	repositoryScene   scene.Interface
	repositoryFile    file.Interface
	cache             cache.Interface
}

// NewService 实例化service
func NewService(repositoryChannel channel.Interface, repositoryScene scene.Interface, repositoryFile file.Interface, cache cache.Interface) *service {
	return &service{
		repositoryChannel: repositoryChannel,
		repositoryScene:   repositoryScene,
		repositoryFile:    repositoryFile,
		cache:             cache,
	}
}
