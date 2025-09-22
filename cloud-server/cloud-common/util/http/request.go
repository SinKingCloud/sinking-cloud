package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client 全局HTTP客户端
var Client = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Request HTTP请求构建器
type Request struct {
	url     string
	method  string
	headers map[string]string
	body    []byte
	client  *http.Client
}

// NewRequest 创建新的HTTP请求构建器
func NewRequest(url, method string) *Request {
	return &Request{
		url:     url,
		method:  method,
		headers: make(map[string]string),
	}
}

// SetHeader 设置单个请求头
func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// SetHeaders 批量设置请求头
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// SetBody 设置请求体（字节数组）
func (r *Request) SetBody(body []byte) *Request {
	r.body = body
	return r
}

// SetJson 设置JSON请求体（自动序列化）
func (r *Request) SetJson(data interface{}) *Request {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.body = nil
		return r
	}
	r.body = jsonData
	return r
}

// SetForm 设置表单请求体（自动构建a=1&b=2&c=3格式）
func (r *Request) SetForm(params map[string]string) *Request {
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}
	r.body = []byte(formData.Encode())
	return r
}

// SetClient 设置自定义HTTP客户端
func (r *Request) SetClient(client *http.Client) *Request {
	r.client = client
	return r
}

// Do 执行请求并返回 header, body, error
func (r *Request) Do() (map[string]string, []byte, error) {
	resp, err := r.DoWithResponse()
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	headers := r.convertHeadersToMap(resp.Header)
	return headers, body, nil
}

// DoWithResponse 执行请求并返回完整的 http.Response 对象
// 注意：调用者需要负责关闭 response.Body
func (r *Request) DoWithResponse() (*http.Response, error) {
	req, err := r.buildRequest()
	if err != nil {
		return nil, err
	}
	client := r.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// buildRequest 构建HTTP请求对象
func (r *Request) buildRequest() (*http.Request, error) {
	var bodyReader io.Reader
	if r.body != nil {
		bodyReader = bytes.NewReader(r.body)
	}
	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// getClient 获取HTTP客户端（优先使用自定义客户端，否则使用全局默认客户端）
func (r *Request) getClient() *http.Client {
	if r.client != nil {
		return r.client
	}
	return Client
}

// convertHeadersToMap 将HTTP响应头转换为map格式
func (r *Request) convertHeadersToMap(headers http.Header) map[string]string {
	result := make(map[string]string, len(headers))
	for k, v := range headers {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
