package rpc

import (
	"encoding/json"
	"service-pay/app/service"
)

type orderService struct {
}

type cleanExpiredUnPayOrderParam struct {
	Day int `json:"day"`
}

// CleanExpiredUnPayOrder 清理过期未支付订单
func (*orderService) CleanExpiredUnPayOrder(params json.RawMessage) (interface{}, error) {
	param := &cleanExpiredUnPayOrderParam{}
	if len(params) > 0 {
		if err := json.Unmarshal(params, param); err != nil {
			return 0, err
		}
	}
	return service.Order.CleanExpiredUnPayOrder(param.Day)
}
