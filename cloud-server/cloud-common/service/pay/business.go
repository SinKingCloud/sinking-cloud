package pay

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CreateOrder 创建订单
func (*Service) CreateOrder(userId int64, name string, describe string, amount float64, payType string, param string, callbackRpcService string, callbackRpcFunc string) (tradeNo string, err error) {
	err = rpc.InvokeFunc(constant.ServicePay, "CreateOrder", &CreateOrderParam{
		UserId:             userId,
		Name:               name,
		Describe:           describe,
		Amount:             amount,
		PayType:            payType,
		Param:              param,
		CallbackRpcService: callbackRpcService,
		CallbackRpcFunc:    callbackRpcFunc,
	}, &tradeNo)
	if err != nil || tradeNo == "" {
		return "", errors.Join(errors.New("创建订单失败"), err)
	}
	return tradeNo, nil
}
