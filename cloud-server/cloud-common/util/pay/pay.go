package pay

import "github.com/go-pay/gopay/alipay"

// NewAliPay 实例化支付宝支付
func NewAliPay(appId, privateKey string, isProd bool) (*alipay.Client, error) {
	return alipay.NewClient(appId, privateKey, isProd)
}
