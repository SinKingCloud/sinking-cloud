package ai

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/net/proxy"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ai struct {
	Proxy      *TCPProxy
	SessionKey string
}

var ins *ai

// NewAi create new ai object
func NewAi() *ai {
	return &ai{}
}

// GetAiIns 获取AI单例
func GetAiIns() *ai {
	if ins == nil {
		ins = NewAi()
	}
	return ins
}

// UseProxy 使用socks5代理
func (a *ai) UseProxy(addr string, port string, user string, pwd string) *ai {
	a.Proxy = &TCPProxy{
		Addr: addr,
		Port: port,
		User: user,
		Pwd:  pwd,
	}
	return a
}

// TCPProxy 代理
type TCPProxy struct {
	Addr string
	Port string
	User string
	Pwd  string
}

// getHttpClient 获取httpClient
func (a *ai) getHttpClient(timeout time.Duration, httpTransport *http.Transport) *http.Client {
	if httpTransport == nil {
		httpTransport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	if a.Proxy != nil && a.Proxy.Addr != "" {
		var a1 *proxy.Auth
		if a.Proxy.User != "" && a.Proxy.Pwd != "" {
			a1 = &proxy.Auth{
				User:     a.Proxy.User,
				Password: a.Proxy.Pwd,
			}
		}
		dialer, e := proxy.SOCKS5("tcp", a.Proxy.Addr+":"+a.Proxy.Port, a1, proxy.Direct)
		if e == nil {
			httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		}
	}
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: httpTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return httpClient
}

// sendHttpRequest 发送http请求
func (a *ai) sendHttpRequest(timeout time.Duration, url string, method string, headers map[string]string, data io.Reader, isGzip bool, noReturn bool) (body []byte, header http.Header, err error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := a.getHttpClient(timeout, nil).Do(req)
	if err != nil {
		return nil, nil, err
	}
	if noReturn {
		return nil, nil, nil
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			e := Body.Close()
			if e != nil {
				return
			}
		}(resp.Body)
	}
	if isGzip {
		r, e := gzip.NewReader(resp.Body)
		if e != nil {
			return nil, nil, e
		}
		b, e := io.ReadAll(r)
		if e != nil {
			return nil, nil, e
		}
		return b, resp.Header, nil
	} else {
		b, e := io.ReadAll(resp.Body)
		if e != nil {
			return nil, nil, e
		}
		return b, resp.Header, nil
	}
}

// SendFormDataRequest 发送formData请求
func (a *ai) SendFormDataRequest(u string, post *map[string]string, file *map[string][]byte) (body []byte, err error) {
	traceId := md5.Sum([]byte(strconv.FormatInt(time.Now().UnixMilli(), 10)))
	header := map[string]string{
		"Referer":      "https://smartapps.cn/AHQ3hWGXDuWutB1dkjwK4g9DrDiDqZFe/0/page-frame.html",
		"User-Agent":   "Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 swan/2.57.0 swan-baiduboxapp/13.19.5.10 baiduboxapp/13.19.5.10 (Baidu; P2 16.1)",
		"X-BD-Traceid": hex.EncodeToString(traceId[:]),
		"sessionKey":   a.GetSessionKey(),
	}
	url := "https://ai.baidu.com" + u
	method := http.MethodGet
	data := bytes.NewBuffer(nil)
	w := multipart.NewWriter(data)
	defer func(w *multipart.Writer) {
		_ = w.Close()
	}(w)
	if post != nil {
		for k, v := range *post {
			_ = w.WriteField(k, v)
		}
		header["Content-Type"] = w.FormDataContentType()
		method = http.MethodPost
	}
	if file != nil {
		for k, v := range *file {
			if createFormFile, e := w.CreateFormFile(k, k); e == nil {
				_, _ = createFormFile.Write(v)
			}
		}
		if method != http.MethodPost {
			header["Content-Type"] = w.FormDataContentType()
			method = http.MethodPost
		}
	}
	respBytes, _, err := a.sendHttpRequest(600*time.Second, url, method, header, data, false, false)
	return respBytes, err
}

// SendNormalRequest 发送普通请求
func (a *ai) SendNormalRequest(u string, post *string) (body []byte, err error) {
	traceId := md5.Sum([]byte(strconv.FormatInt(time.Now().UnixMilli(), 10)))
	header := map[string]string{
		"Referer":      "https://smartapps.cn/AHQ3hWGXDuWutB1dkjwK4g9DrDiDqZFe/0/page-frame.html",
		"User-Agent":   "Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 swan/2.57.0 swan-baiduboxapp/13.19.5.10 baiduboxapp/13.19.5.10 (Baidu; P2 16.1)",
		"X-BD-Traceid": hex.EncodeToString(traceId[:]),
		"sessionKey":   a.GetSessionKey(),
	}
	url := "https://ai.baidu.com" + u
	method := http.MethodGet
	data := strings.NewReader(*post)
	if post != nil {
		method = http.MethodPost
	}
	respBytes, _, err := a.sendHttpRequest(600*time.Second, url, method, header, data, false, false)
	return respBytes, err
}

// GetSessionKey 获取sessionKey
func (a *ai) GetSessionKey() string {
	if a.SessionKey == "" {
		url := "https://ai.baidu.com/weapp/login?code=4908f2e82a54121da3a823a5a8b9e860Nn"
		header := map[string]string{
			"Referer":    "https://smartapps.cn/AHQ3hWGXDuWutB1dkjwK4g9DrDiDqZFe/0/page-frame.html",
			"User-Agent": "Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 swan/2.57.0 swan-baiduboxapp/13.19.5.10 baiduboxapp/13.19.5.10 (Baidu; P2 16.1)",
		}
		body, _, err := a.sendHttpRequest(60*time.Second, url, http.MethodGet, header, nil, false, false)
		if err == nil {
			type Response struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
				Data struct {
					SessionKey string `json:"sessionKey"`
				} `json:"data"`
			}
			var resp2 *Response
			err = json.Unmarshal(body, &resp2)
			if err != nil || resp2 == nil || resp2.Code != 0 {
				a.SessionKey = resp2.Msg
			} else {
				a.SessionKey = resp2.Data.SessionKey
			}
		} else {
			a.SessionKey = err.Error()
		}
	}
	return a.SessionKey
}
