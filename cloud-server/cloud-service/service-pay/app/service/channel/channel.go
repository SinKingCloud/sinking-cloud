package channel

import (
	"net/http"
	"service-pay/app/repository/channel"
	"service-pay/app/repository/order"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
	"github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

// Service service接口
type Service interface {
	Create(name string, Type string, plugin int, status int, config string) error                                                                                                      //创建数据
	Select(where *channel.SelectChannel, page *page.Query) (list []*channel.Channel, total int64, err error)                                                                           //查询数据
	DeleteById(id int64) error                                                                                                                                                         //通过ID删除数据
	UpdateById(id int64, data *channel.UpdateChannel) error                                                                                                                            //通过ID更新数据
	FindById(id int64, refresh bool) (*Channel, error)                                                                                                                                 //通过ID查询数据
	GetConfig(refresh bool) *Config                                                                                                                                                    //获取支付配置
	SetConfig(data *Config) error                                                                                                                                                      //更新支付配置
	GetPayResultPageUrl(code int) string                                                                                                                                               //获取支付结果页地址
	GetPayUrl(channelId int64, orderId int64, name string, money float64) (string, error)                                                                                              //获取支付跳转链接
	ValidateCallback(channelId int64, req *http.Request, successFun func(data *pay.CallbackData) error) (code int, header http.Header, body string, data *pay.CallbackData, err error) //验证支付回调
	GetPayChannel(payType string) (*Channel, error)                                                                                                                                    //获取支付通道ID
	GetIdNameMap(refresh bool) (map[int64]string, error)                                                                                                                               //获取id name映射
}

// Session服务
type service struct {
	repositoryChannel channel.Interface
	orderRepository   order.Interface
	configService     *config.Service
	consul            *client.Client
	cache             cache.Interface
}

// NewService 实例化service
func NewService(repositoryChannel channel.Interface, orderRepository order.Interface, configService *config.Service, consul *client.Client, cache cache.Interface) *service {
	return &service{
		repositoryChannel: repositoryChannel,
		orderRepository:   orderRepository,
		configService:     configService,
		consul:            consul,
		cache:             cache,
	}
}
