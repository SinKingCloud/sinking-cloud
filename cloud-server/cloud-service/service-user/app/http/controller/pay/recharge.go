package pay

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/order"
	"strconv"
)

func (ControllerPay) Recharge(c *sinking_server.Context) {
	type Form struct {
		Money int `json:"money" default:"" validate:"required,numeric,min=1,max=99999" label:"充值金额"`
		Type  int `json:"type" default:"" validate:"oneof=0 1" label:"充值方式"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	min, _ := strconv.Atoi(configs[config.PayMinMoney])
	if form.Money < min {
		c.Error("系统限制最低充值" + strconv.Itoa(min) + "元")
		return
	}
	//创建订单
	tradeNo := str.GetSnowWorkIns().GetId()
	err := order.GetIns().Create(&model.WebUserOrder{
		UserId:    c.GetWebUserInfo().Id,
		PayType:   form.Type,
		OrderType: 0,
		Money:     float64(form.Money),
		TradeNo:   strconv.FormatInt(tradeNo, 10),
		Name:      "账户余额充值",
		Status:    0,
	})
	if err == nil {
		c.SuccessWithData("创建订单成功", tradeNo)
	} else {
		c.Error("创建订单失败")
	}
}
