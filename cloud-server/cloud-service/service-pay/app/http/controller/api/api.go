package api

import (
	"net/http"
	"service-pay/app/enum/order_status"
	"service-pay/app/enum/pay_type"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
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
		pay_type.Balance: config.Balance > 0,
		pay_type.Alipay:  config.Alipay > 0,
		pay_type.Wxpay:   config.Wxpay > 0,
	})
}

// Submit 获取支付链接
func (Controller) Submit(c *context.Context) {
	var form struct {
		TradeNo string `json:"trade_no" default:"" validate:"required,numeric" label:"订单号"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Order.FindByTradeNo(form.TradeNo)
	if err != nil {
		c.Error(err.Error())
		return
	}
	if info.Status != order_status.UnPay {
		c.Error("该订单已支付或已关闭")
		return
	}
	link, err := service.Channel.GetPayUrl(info.ChannelId, info.TradeNo, info.Name, info.Amount)
	if err != nil {
		c.Error(err.Error())
		return
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
	code, header, body, _, err := service.Channel.ValidateCallback(id, c.Request, func(data *pay.CallbackData) bool {
		err := service.Order.Callback(data)
		if err == nil {
			return true
		}
		return false
	})
	if method == "return" {
		if err == nil {
			c.SetHeader("content-type", "text/html;charset=utf-8;")
			c.Data(http.StatusOK, []byte("<script>window.location.href=\"/?s#/pay\"</script>"))
		} else {
			c.String(http.StatusOK, "充值失败,请联系管理员")
		}
	} else {
		for k, v := range header {
			c.SetHeader(k, v[0])
		}
		c.String(code, body)
	}
}
