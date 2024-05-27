package callback

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"net/http"
)

// Alipay 支付宝支付
func (ControllerCallback) Alipay(c *sinking_server.Context) {
	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.PayGroup, false)
	ok, err := alipay.VerifySign(configs[config.PayAlipayPublicKey], notifyReq)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if !ok {
		c.String(http.StatusOK, "fail")
		return
	}
	bm := make(gopay.BodyMap)
	err = notifyReq.Unmarshal(&bm)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	//执行业务逻辑
	if service2.Order.CallBack(bm.GetString("out_trade_no"),
		bm.GetString("trade_no")) {
		c.String(http.StatusOK, "success")
	} else {
		c.String(http.StatusOK, "fail")
	}
}
