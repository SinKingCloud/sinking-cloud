package pay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// EpayGateway Epay支付网关实现（易支付）
type EpayGateway struct {
	pid       string
	key       string
	submitURL string
	mapiURL   string
	apiURL    string
	signType  string
}

// NewEpayGateway 创建Epay网关实例
// pid: 商户ID
// key: 商户密钥
// url: 易支付接口地址
// channel: 渠道标识（1或2）
func NewEpayGateway(pid, key, apiUrl string) (*EpayGateway, error) {
	if pid == "" || key == "" || apiUrl == "" {
		return nil, errors.New("初始化易支付客户端失败: 参数不能为空")
	}
	if !strings.HasSuffix(apiUrl, "/") {
		apiUrl = apiUrl + "/"
	}
	return &EpayGateway{
		pid:       pid,
		key:       key,
		submitURL: apiUrl + "submit.php",
		mapiURL:   apiUrl + "mapi.php",
		apiURL:    apiUrl + "api.php",
		signType:  "MD5",
	}, nil
}

// Submit 提交支付（实现PaymentGateway接口）
func (e *EpayGateway) Submit(order *Order) (string, error) {
	if order == nil || order.Id <= 0 {
		return "", errors.New("订单ID错误")
	}
	// 构建支付参数
	params := map[string]string{
		"pid":          e.pid,
		"type":         order.PayType,
		"notify_url":   order.NotifyUrl,
		"return_url":   order.ReturnUrl,
		"out_trade_no": strconv.FormatInt(order.Id, 10),
		"name":         order.Name,
		"money":        fmt.Sprintf("%.2f", order.Money),
	}
	// 生成支付链接
	return e.getPayLink(params), nil
}

// ParseCallback 解析回调参数（实现 Interface 接口）
func (e *EpayGateway) ParseCallback(req *http.Request) (*CallbackData, error) {
	params := make(map[string]string)
	if req.Method == "GET" {
		for k, v := range req.URL.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	} else {
		_ = req.ParseForm()
		for k, v := range req.PostForm {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}

	// 验证签名
	if len(params) == 0 {
		return nil, errors.New("参数为空")
	}
	sign := e.getSign(params)
	if sign != params["sign"] {
		return nil, errors.New("签名验证失败")
	}

	// 解析金额
	money, err := strconv.ParseFloat(params["money"], 64)
	if err != nil {
		money = 0
	}
	orderId, err := strconv.ParseInt(strings.TrimSpace(params["out_trade_no"]), 10, 64)
	if err != nil || orderId <= 0 {
		return nil, errors.New("订单ID错误")
	}

	return &CallbackData{
		Id:         orderId,
		OutTradeNo: params["trade_no"],
		Status:     params["trade_status"] == "TRADE_SUCCESS",
		Money:      money,
		PayTime:    params["time"],
	}, nil
}

// SuccessResponse 获取成功响应状态码、响应头和响应内容
func (e *EpayGateway) SuccessResponse() (int, http.Header, string) {
	return http.StatusOK, http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}, "success"
}

// FailResponse 获取失败响应状态码、响应头和响应内容
func (e *EpayGateway) FailResponse() (int, http.Header, string) {
	return http.StatusOK, http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}}, "fail"
}

// QueryOrder 查询订单状态（实现PaymentGateway接口）
func (e *EpayGateway) QueryOrder(orderId int64) (*OrderStatus, error) {
	if orderId <= 0 {
		return nil, errors.New("订单ID错误")
	}
	result := e.queryOrderAPI(strconv.FormatInt(orderId, 10))
	if result == nil {
		return nil, errors.New("查询失败")
	}
	status := false
	if val, ok := result["status"].(float64); ok {
		status = val == 1
	}
	money := 0.0
	if val, ok := result["money"].(float64); ok {
		money = val
	}
	outTradeNo := ""
	if val, ok := result["trade_no"].(string); ok {
		outTradeNo = val
	}
	return &OrderStatus{
		Status:     status,
		OutTradeNo: outTradeNo,
		Money:      money,
	}, nil
}

// getPayLink 生成支付链接
func (e *EpayGateway) getPayLink(params map[string]string) string {
	builtParams := e.buildRequestParams(params)
	urlValues := make(map[string][]string)
	for k, v := range builtParams {
		urlValues[k] = []string{v}
	}
	return e.submitURL + "?" + url.Values(urlValues).Encode()
}

// getPagePayHTML 生成支付表单HTML（页面跳转方式）
func (e *EpayGateway) getPagePayHTML(params map[string]string, button string) string {
	builtParams := e.buildRequestParams(params)
	html := fmt.Sprintf(`<form id="dopay" action="%s" method="post">`, e.submitURL)
	for k, v := range builtParams {
		html += fmt.Sprintf(`<input type="hidden" name="%s" value="%s"/>`, k, v)
	}
	html += fmt.Sprintf(`<input type="submit" value="%s"></form><script>document.getElementById("dopay").submit();</script>`, button)
	return html
}

// queryOrderAPI 调用易支付API查询订单
func (e *EpayGateway) queryOrderAPI(orderId string) map[string]interface{} {
	queryURL := fmt.Sprintf("%s?act=order&pid=%s&key=%s&out_trade_no=%s", e.apiURL, e.pid, e.key, orderId)
	response, err := http.Get(queryURL)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil
	}
	return result
}

// buildRequestParams 构建请求参数（添加签名）
func (e *EpayGateway) buildRequestParams(params map[string]string) map[string]string {
	sign := e.getSign(params)
	params["sign"] = sign
	params["sign_type"] = e.signType
	return params
}

// getSign 计算签名
func (e *EpayGateway) getSign(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	signStr := ""
	for _, k := range keys {
		v := params[k]
		if k != "sign" && k != "sign_type" && v != "" {
			signStr += fmt.Sprintf("%s=%s&", k, v)
		}
	}
	signStr = strings.TrimSuffix(signStr, "&")
	signStr += e.key
	return str.NewStringTool(signStr).Md5()
}
