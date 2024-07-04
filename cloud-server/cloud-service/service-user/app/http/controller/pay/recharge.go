package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/order"
	"strconv"
)

func (ControllerPay) Recharge(c *sinking_server.Context) {
	type Form struct {
		Money int `form:"money" json:"money" default:"" validate:"min=1,max=99999" label:"充值金额"`
		Type  int `form:"type" json:"type" default:"" validate:"oneof=0" label:"充值方式"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.PayGroup, false)
	min, _ := strconv.Atoi(configs[config.PayMinMoney])
	if form.Money < min {
		c.Error("系统限制最低充值" + strconv.Itoa(min) + "元")
		return
	}
	//创建订单
	tradeNo := strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10)
	err := service2.Order.Create(&model.WebUserOrder{
		UserId:    c.GetWebUserInfo().Id,
		PayType:   form.Type,
		OrderType: order.OrdersTypeRecharge,
		Money:     float64(form.Money),
		TradeNo:   tradeNo,
		Name:      "账户余额充值",
		Status:    order.StatusUnPay,
	})
	if err == nil {
		c.SuccessWithData("创建订单成功", tradeNo)
	} else {
		c.Error("创建订单失败")
	}
}
