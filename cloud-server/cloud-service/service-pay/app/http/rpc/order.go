package rpc

import (
	"encoding/json"
	"service-pay/app/enum/order_type"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
)

type orderService struct {
}

// CleanExpiredUnPayOrder 清理30天前过期未支付订单
func (*orderService) CleanExpiredUnPayOrder(_ json.RawMessage) (interface{}, error) {
	return service.Order.CleanExpiredUnPayOrder(30)
}

// CreateOrder 生成订单
func (*orderService) CreateOrder(params json.RawMessage) (interface{}, error) {
	var param *pay.CreateOrderParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return "", err
	}
	tradeNo, err := service.Order.Create(param.UserId,
		param.Name,
		order_type.Buy,
		param.PayType,
		param.Amount,
		param.Param,
		param.CallbackRpcService,
		param.CallbackRpcFunc, param.Describe)
	if err != nil {
		return "", err
	}
	return tradeNo, nil
}
