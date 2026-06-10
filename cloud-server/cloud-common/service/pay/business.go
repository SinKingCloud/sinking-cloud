package pay

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CreateOrder 创建订单
func (*Service) CreateOrder(userId int64, name string, describe string, amount float64, payType string, param string, callbackRpcService string, callbackRpcFunc string) (orderId int64, err error) {
	err = rpc.InvokeFunc(constant.ServicePay, "CreateOrder", &CreateOrderParam{
		UserId:             userId,
		Name:               name,
		Describe:           describe,
		Amount:             amount,
		PayType:            payType,
		Param:              param,
		CallbackRpcService: callbackRpcService,
		CallbackRpcFunc:    callbackRpcFunc,
	}, &orderId)
	if err != nil || orderId <= 0 {
		return 0, errors.Join(errors.New("创建订单失败"), err)
	}
	return orderId, nil
}

// CleanExpiredUnPayOrder 清理过期未支付订单
func (*Service) CleanExpiredUnPayOrder(day int) (successNum int64, err error) {
	err = rpc.InvokeFunc(constant.ServicePay, "CleanExpiredUnPayOrder", day, &successNum)
	if err != nil {
		return 0, errors.Join(errors.New("清理过期未支付订单失败"), err)
	}
	return successNum, nil
}
