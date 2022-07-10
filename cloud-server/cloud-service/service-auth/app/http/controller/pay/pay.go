package pay

import (
	"context"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/pay"
	"github.com/go-pay/gopay"
)

type ControllerPay struct {
}

// Alipay 支付宝支付
func (ControllerPay) Alipay(tradeNo string, money string, name string) (url string, err error) {
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.PayGroup, false)
	configs2 := service.GetSystemConfigServiceIns().GetGroupByCache(config.SystemGroup, false)
	client, err := pay.NewAliPay(configs[config.PayAlipayAppid], configs[config.PayAlipayPrivateKey], true)
	if err != nil {
		err = errors.New("初始化支付请求失败")
		return
	}
	ctx := context.Background()
	client.SetReturnUrl(configs2[config.SystemUrl] + "/pay/recharge").
		SetNotifyUrl(configs2[config.SystemGateway] + "/auth/callback/alipay")
	bm := make(gopay.BodyMap)
	bm.Set("subject", name).
		Set("out_trade_no", tradeNo).
		Set("total_amount", money).
		Set("product_code", "FAST_INSTANT_TRADE_PAY").
		Set("timeout_express", "1h")
	url, err = client.TradePagePay(ctx, bm)
	if err != nil {
		err = errors.New("请求支付接口失败")
		return
	}
	ctx.Done()
	return
}
