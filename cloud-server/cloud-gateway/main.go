package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	go func() {
		wp, err := NewProxy("ws://127.0.0.1:1004/test/1", func(r *http.Request) error {
			// 权限验证
			r.Header.Set("Cookie", "----")
			// 伪装来源
			r.Header.Set("Origin", "http://82.157.123.54:9010")
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
		http.HandleFunc("/test", wp.Proxy)
		http.ListenAndServe(":9696", nil)
	}()
	bootstrap.Load()
	app.Run()
}

const (
	WsScheme  = "ws"
	WssScheme = "wss"
)

var ErrFormatAddr = errors.New("remote websockets addr format error")

func NewProxy(addr string, beforeCallback func(r *http.Request) error) (*WebsocketProxy, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, ErrFormatAddr
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, ErrFormatAddr
	}
	if u.Scheme != WsScheme && u.Scheme != WssScheme {
		return nil, ErrFormatAddr
	}
	wp := &WebsocketProxy{
		scheme:          u.Scheme,
		remoteAddr:      fmt.Sprintf("%s:%s", host, port),
		defaultPath:     u.Path,
		beforeHandshake: beforeCallback,
	}
	if u.Scheme == WssScheme {
		wp.tlsc = &tls.Config{InsecureSkipVerify: true} // 不验证证书
	}
	return wp, nil
}

type WebsocketProxy struct {
	scheme          string // ws, wss
	remoteAddr      string // 目标地址: host:port
	defaultPath     string // path地址
	tlsc            *tls.Config
	beforeHandshake func(r *http.Request) error // 发送握手之前回调
}

func (wp *WebsocketProxy) Proxy(writer http.ResponseWriter, request *http.Request) {
	// 判断是否是websocket请求
	if strings.ToLower(request.Header.Get("Connection")) != "upgrade" ||
		strings.ToLower(request.Header.Get("Upgrade")) != "websocket" {
		_, _ = writer.Write([]byte(`Must be a websocket request`))
		return
	}
	// 劫持连接
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		return
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return
	}
	defer conn.Close()
	// 克隆请求，设置目标地址路径
	req := request.Clone(context.TODO())
	req.URL.Path, req.URL.RawPath, req.RequestURI = wp.defaultPath, wp.defaultPath, wp.defaultPath
	req.Host = wp.remoteAddr
	// 握手之前回调
	if wp.beforeHandshake != nil {
		// 增加头部，权限认证 + 伪装来源
		err = wp.beforeHandshake(req)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
	}
	// 判断协议，选择拨号流程
	var remoteConn net.Conn
	switch wp.scheme {
	case WsScheme:
		remoteConn, err = net.Dial("tcp", wp.remoteAddr)
	case WssScheme:
		remoteConn, err = tls.Dial("tcp", wp.remoteAddr, wp.tlsc)
	}
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	defer remoteConn.Close()
	// request请求转换为字节流
	b, _ := httputil.DumpRequest(req, false)
	// 向目标websocket服务发送握手包
	remoteConn.Write(b)
	// 流量透传
	errChan := make(chan error, 2)
	copyConn := func(a, b net.Conn) {
		_, err := io.Copy(a, b)
		errChan <- err
	}
	go copyConn(conn, remoteConn) // response
	go copyConn(remoteConn, conn) // request
	select {
	case err = <-errChan:
		if err != nil {
			log.Println(err)
		}
	}
}
