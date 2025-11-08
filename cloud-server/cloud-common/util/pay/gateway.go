package pay

import "net/http"

// Interface 支付网关接口
type Interface interface {
	// Submit 提交支付，返回支付链接或表单
	Submit(order *Order) (string, error)

	// ParseCallback 解析回调参数
	ParseCallback(req *http.Request) (*CallbackData, error)

	// QueryOrder 查询订单状态
	QueryOrder(tradeNo string) (*OrderStatus, error)
}

// Order 支付订单信息
type Order struct {
	TradeNo   string  // 商户订单号
	Name      string  // 商品名称
	Money     float64 // 金额
	NotifyUrl string  // 异步通知地址
	ReturnUrl string  // 同步返回地址
	PayType   string  // 支付类型
}

// CallbackData 回调数据
type CallbackData struct {
	TradeNo    string  // 商户订单号
	OutTradeNo string  // 第三方订单号
	Money      float64 // 支付金额
	Status     string  // 支付状态
	PayTime    string  // 支付时间
}

// OrderStatus 订单状态
type OrderStatus struct {
	Status     bool    // 是否支付成功
	Money      float64 // 支付金额
	OutTradeNo string  // 第三方订单号
}
