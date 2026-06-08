package rpc

import (
	"encoding/json"
	"service-pay/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/pay"
)

type balanceService struct {
}

// RechargeOrderNotify 充值订单回调
func (*balanceService) RechargeOrderNotify(params json.RawMessage) (interface{}, error) {
	var param *pay.NotifyParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return false, err
	}
	err = service.Balance.RechargeOrderNotify(param)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetBalanceByUserIds 批量查询用户余额
func (*balanceService) GetBalanceByUserIds(params json.RawMessage) (interface{}, error) {
	var userIds []int64
	if err := json.Unmarshal(params, &userIds); err != nil {
		return nil, err
	}
	return service.Balance.GetByUserIds(userIds)
}
