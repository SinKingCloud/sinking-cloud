package pay

import "net/http"

// Interface 支付网关接口
type Interface interface {
	// Submit 提交支付，返回支付链接或表单
	Submit(order *Order) (string, error)

	// ParseCallback 解析回调参数
	ParseCallback(req *http.Request) (*CallbackData, error)

	// SuccessResponse 获取成功响应状态码、响应头和响应内容
	SuccessResponse() (int, http.Header, string)

	// FailResponse 获取失败响应状态码、响应头和响应内容
	FailResponse() (int, http.Header, string)

	// QueryOrder 查询订单状态
	QueryOrder(orderId int64) (*OrderStatus, error)
}

// Order 支付订单信息
type Order struct {
	Id        int64   // 系统订单ID
	Name      string  // 商品名称
	Money     float64 // 金额
	NotifyUrl string  // 异步通知地址
	ReturnUrl string  // 同步返回地址
	PayType   string  // 支付类型
}

// CallbackData 回调数据
type CallbackData struct {
	Id      int64   // 系统订单ID
	TradeNo string  // 第三方订单号
	Money   float64 // 支付金额
	Status  bool    // 是否支付成功
	PayTime string  // 支付时间
}

// OrderStatus 订单状态
type OrderStatus struct {
	Status  bool    // 是否支付成功
	Money   float64 // 支付金额
	TradeNo string  // 第三方订单号
}
