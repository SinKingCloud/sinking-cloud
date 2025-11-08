package pay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

// AlipayGateway 支付宝官方SDK网关
type AlipayGateway struct {
	client          *alipay.Client
	appId           string
	aliPayPublicKey string // 保存支付宝公钥用于验签
}

// NewAlipayGateway 创建支付宝官方网关
// appId: 应用ID
// privateKey: 应用私钥（PKCS1或PKCS8格式的字符串）
// aliPayPublicKey: 支付宝公钥
func NewAlipayGateway(appId, privateKey, aliPayPublicKey string) (*AlipayGateway, error) {
	client, err := alipay.NewClient(appId, privateKey, true)
	if err != nil {
		return nil, fmt.Errorf("初始化支付宝客户端失败: %w", err)
	}
	// 关闭调试日志
	client.DebugSwitch = gopay.DebugOff
	return &AlipayGateway{
		client:          client,
		appId:           appId,
		aliPayPublicKey: aliPayPublicKey,
	}, nil
}

// Submit 提交支付（电脑网站支付）
func (a *AlipayGateway) Submit(order *Order) (string, error) {
	ctx := context.Background()
	// 构建请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", order.Name)
	bm.Set("out_trade_no", order.TradeNo)
	bm.Set("total_amount", fmt.Sprintf("%.2f", order.Money))
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	// 设置回调地址
	if order.NotifyUrl != "" {
		bm.Set("notify_url", order.NotifyUrl)
	}
	if order.ReturnUrl != "" {
		bm.Set("return_url", order.ReturnUrl)
	}
	// 调用支付宝电脑网站支付API
	payUrl, err := a.client.TradePagePay(ctx, bm)
	if err != nil {
		return "", fmt.Errorf("支付宝下单失败: %w", err)
	}
	return payUrl, nil
}

// ParseCallback 解析回调参数（实现 Interface 接口）
func (a *AlipayGateway) ParseCallback(req *http.Request) (*CallbackData, error) {
	notifyReq, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return nil, fmt.Errorf("解析回调参数失败: %w", err)
	}
	ok, err := alipay.VerifySign(a.aliPayPublicKey, notifyReq)
	if err != nil {
		return nil, fmt.Errorf("验签失败: %w", err)
	}
	if !ok {
		return nil, errors.New("签名验证失败")
	}
	bm := make(gopay.BodyMap)
	err = notifyReq.Unmarshal(&bm)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败: %w", err)
	}
	money, err := strconv.ParseFloat(bm.GetString("total_amount"), 64)
	if err != nil {
		return nil, fmt.Errorf("金额解析失败: %w", err)
	}
	return &CallbackData{
		TradeNo:    bm.GetString("out_trade_no"), // 商户订单号
		OutTradeNo: bm.GetString("trade_no"),     // 支付宝交易号
		Status:     "TRADE_SUCCESS",              // TRADE_SUCCESS 或 TRADE_FINISHED
		Money:      money,
		PayTime:    bm.GetString("gmt_payment"), // 支付时间
	}, nil
}

// QueryOrder 查询订单状态
func (a *AlipayGateway) QueryOrder(tradeNo string) (*OrderStatus, error) {
	ctx := context.Background()
	// 构建请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", tradeNo)
	// 调用支付宝查询订单API
	rsp, err := a.client.TradeQuery(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	// 检查返回状态 (10000表示成功)
	if rsp.Response.Code != "10000" {
		return nil, fmt.Errorf("查询订单失败: %s - %s", rsp.Response.Code, rsp.Response.Msg)
	}
	// 解析金额
	money, _ := strconv.ParseFloat(rsp.Response.TotalAmount, 64)
	// 判断订单状态
	// WAIT_BUYER_PAY: 交易创建，等待买家付款
	// TRADE_CLOSED: 未付款交易超时关闭，或支付完成后全额退款
	// TRADE_SUCCESS: 交易支付成功
	// TRADE_FINISHED: 交易结束，不可退款
	isSuccess := rsp.Response.TradeStatus == "TRADE_SUCCESS" || rsp.Response.TradeStatus == "TRADE_FINISHED"
	return &OrderStatus{
		Status:     isSuccess,
		OutTradeNo: rsp.Response.TradeNo,
		Money:      money,
	}, nil
}
