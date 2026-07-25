package api

import (
	"service-pay/app/enum/order_status"
	"service-pay/app/service"
	"strconv"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Submit 获取支付链接
func (Controller) Submit(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"0" validate:"numeric,min=1" label:"订单ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Order.FindById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	if info.Status != order_status.UnPay {
		c.Error("该订单已支付或已关闭,无需支付")
		return
	}
	expiredTime := time.Time(info.ExpiredTime)
	if !expiredTime.IsZero() && !expiredTime.After(time.Now()) {
		c.Error("该订单已过期,请重新下单")
		return
	}
	var link string
	if info.Amount <= 0 {
		if err = service.Order.Callback(info.Id, strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10), 0, str.DateTime(time.Now())); err != nil {
			c.Error(err.Error())
			return
		}
		link = service.Channel.GetPayResultPageUrl(200)
	} else {
		link, err = service.Channel.GetPayUrl(info.ChannelId, info.Id, info.Name, info.Amount)
		if err != nil {
			c.Error(err.Error())
			return
		}
	}
	c.SuccessWithData("获取支付信息成功", link)
}
