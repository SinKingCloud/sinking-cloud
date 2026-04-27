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
	var param *pay.NotifyInfo
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
