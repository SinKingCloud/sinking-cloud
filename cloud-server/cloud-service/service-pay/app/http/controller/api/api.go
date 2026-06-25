package api

import (
	"errors"
	"net/http"
	"service-pay/app/enum/order_status"
	"service-pay/app/enum/pay_type"
	"service-pay/app/service"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type Controller struct {
}

// Config 支付配置
func (Controller) Config(c *context.Context) {
	config := service.Channel.GetConfig(false)
	c.SuccessWithData("获取数据成功", sinking_web.H{
		"min":            config.Min,
		"max":            config.Max,
		"expire":         config.Expire,
		pay_type.Balance: config.Balance > 0,
		pay_type.Alipay:  config.Alipay > 0,
		pay_type.Wxpay:   config.Wxpay > 0,
	})
}

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

// Callback 支付回调
func (Controller) Callback(c *context.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	method := c.Param("method")
	if id <= 0 || method == "" {
		c.Error("参数错误")
		return
	}
	code, header, body, _, err := service.Channel.ValidateCallback(id, c.Request, func(data *pay.CallbackData) error {
		if data == nil {
			return errors.New("回调信息不能为空")
		}
		if data.Id <= 0 {
			return errors.New("订单ID错误")
		}
		payTime := str.DateTime(time.Now())
		payTimeStr := strings.TrimSpace(data.PayTime)
		if payTimeStr != "" {
			if t, err := time.ParseInLocation("2006-01-02 15:04:05", payTimeStr, time.Local); err == nil {
				payTime = str.DateTime(t)
			}
		}
		return service.Order.Callback(data.Id, strings.TrimSpace(data.TradeNo), data.Money, payTime)
	})
	if method == "return" {
		res := 200
		if err != nil {
			res = 500
		}
		c.Redirect(http.StatusFound, service.Channel.GetPayResultPageUrl(res))
	} else {
		for k, v := range header {
			c.SetHeader(k, v[0])
		}
		c.String(code, body)
	}
}
