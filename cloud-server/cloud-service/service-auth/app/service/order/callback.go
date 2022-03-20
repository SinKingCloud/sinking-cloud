package order

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/service/pay_log"
)

func (s *service) CallBack(tradeNo string, outTradeNo string) bool {
	order, err := s.FindByTradeNo(tradeNo)
	if err != nil {
		return false
	}
	if order.Status != 0 {
		return true
	}
	money := fmt.Sprintf("%.2f", order.Money)
	res := false
	switch order.OrderType {
	case 0:
		//余额充值
		e := service2.GetWebUserServiceIns().Money(order.UserId, money, 1)
		if e == nil {
			res = true
			_ = pay_log.GetIns().Create(&model.WebUserPayLog{
				UserId:  order.UserId,
				Type:    0,
				LogType: 0,
				Title:   "在线充值",
				Content: "在线充值余额" + money + "元",
				Money:   order.Money,
			})
			service2.GetWebUserServiceIns().ClearCache(order.UserId)
		}
		break
	}
	//更改订单状态
	if res {
		_ = s.UpdateById(order.Id, map[string]interface{}{"status": "1", "out_trade_no": outTradeNo})
	}
	return res
}
