package rpc

import (
	"encoding/json"
	"service-pay/app/service"
)

type orderService struct {
}

// CleanExpiredUnPayOrder 清理30天前过期未支付订单
func (*orderService) CleanExpiredUnPayOrder(_ json.RawMessage) (interface{}, error) {
	return service.Order.CleanExpiredUnPayOrder(30)
}
