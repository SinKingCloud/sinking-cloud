package file

import (
	"context"
	"io"
	"service-driver/app/model"
	"service-driver/app/repository/channel"
	"service-driver/app/repository/file"
	"service-driver/app/repository/scene"
	channel2 "service-driver/app/service/channel"
	scene2 "service-driver/app/service/scene"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	SimpleUpload(ctx context.Context, sceneId int64, body io.Reader, filename string, size int64) (string, error)                                //简单上传
	Select(where *file.SelectFile, orderByField string, orderByType string, page int, pageSize int) (list []*model.File, total int64, err error) //查询数据
	GetUrlByHash(hash string) (string, error)                                                                                                    //根据hash获取文件公开地址
	UpdateByIds(ids []int64, data *file.UpdateFile) error                                                                                        //批量更新数据
	DeleteByIds(ids []int64, clear bool) error                                                                                                   //批量删除数据
}

// Session服务
type service struct {
	repositoryChannel channel.Interface
	repositoryScene   scene.Interface
	repositoryFile    file.Interface
	channelService    channel2.Service
	sceneService      scene2.Service
	cache             cache.Interface
}

// NewService 实例化service
func NewService(repositoryChannel channel.Interface, repositoryScene scene.Interface, repositoryFile file.Interface, channelService channel2.Service, sceneService scene2.Service, cache cache.Interface) *service {
	return &service{
		repositoryChannel: repositoryChannel,
		repositoryScene:   repositoryScene,
		repositoryFile:    repositoryFile,
		channelService:    channelService,
		sceneService:      sceneService,
		cache:             cache,
	}
}
