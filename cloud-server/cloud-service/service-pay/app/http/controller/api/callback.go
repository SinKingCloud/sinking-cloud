package api

import (
	"errors"
	"net/http"
	"service-pay/app/service"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

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
		c.String(code, "%s", body)
	}
}
