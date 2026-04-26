package channel

import (
	"service-pay/app/repository/channel"
	"service-pay/app/repository/order"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	Create(name string, types int, status int, config string) error                                                                                         //创建数据
	Select(where *channel.SelectChannel, orderByField string, orderByType string, page int, pageSize int) (list []*channel.Channel, total int64, err error) //查询数据
	DeleteById(id int64) error                                                                                                                              //通过ID删除数据
	UpdateById(id int64, data *channel.UpdateChannel) error                                                                                                 //通过ID更新数据
	FindById(id int64, refresh bool) (*Channel, error)                                                                                                      //通过ID查询数据
	GetConfig(refresh bool) *Config                                                                                                                         //获取支付配置
	SetConfig(data *Config) error                                                                                                                           //更新支付配置
}

// Session服务
type service struct {
	repositoryChannel channel.Interface
	orderRepository   order.Interface
	configService     *config.Service
	cache             cache.Interface
}

// NewService 实例化service
func NewService(repositoryChannel channel.Interface, orderRepository order.Interface, configService *config.Service, cache cache.Interface) *service {
	return &service{
		repositoryChannel: repositoryChannel,
		orderRepository:   orderRepository,
		configService:     configService,
		cache:             cache,
	}
}
