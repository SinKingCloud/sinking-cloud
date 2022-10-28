package ai

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"time"
)

type ai struct {
	Proxy *TCPProxy
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
