package order

import (
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/pay_log"
)

func (s *Service) CallBack(tradeNo string, outTradeNo string) bool {
	key := constant.LockChangeOrderStatus + tradeNo
	if !util.Cache.Lock(key, constant.LockTimeChangeOrderStatus) {
		return false
	}
	defer util.Cache.UnLock(key)
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
	case OrdersTypeRecharge:
		var giftMoney float64
		configs := service2.GetSystemConfigServiceIns().GetGroupByCache(config.PayGroup, false)
		if configs[config.PayGift] == "1" && configs[config.PayGiftMoney] != "" {
			var gift map[string]float64
			e := json.Unmarshal([]byte(configs[config.PayGiftMoney]), &gift)
			if e == nil && gift != nil && gift[fmt.Sprintf("%.2f", order.Money)] > 0 {
				giftMoney = gift[fmt.Sprintf("%.2f", order.Money)]
				money = fmt.Sprintf("%.2f", order.Money+giftMoney)
			}
		}
		//余额充值
		e := service2.GetWebUserServiceIns().Money(order.UserId, money, 1)
		if e == nil {
			res = true
			gift := ""
			if giftMoney > 0 {
				money = fmt.Sprintf("%.2f", order.Money)
				gift = ",平台赠送" + fmt.Sprintf("%.2f", giftMoney) + "元"
			}
			_ = pay_log.GetIns().Create(&model.WebUserPayLog{
				UserId:  order.UserId,
				Type:    pay_log.TypeAdd,
				LogType: pay_log.LogTypeRecharge,
				Title:   "在线充值",
				Content: "在线充值余额" + money + "元" + gift,
				Money:   order.Money + giftMoney,
			})
			service2.GetWebUserServiceIns().ClearCache(order.UserId)
		}
		break
	}
	//更改订单状态
	if res {
		_ = s.UpdateById(order.Id, map[string]interface{}{"status": StatusPayed, "out_trade_no": outTradeNo})
	}
	return res
}
