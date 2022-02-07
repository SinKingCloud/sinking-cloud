package origin

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

//请求client
var httpClient *http.Client

// GetClient 获取http client
func GetClient() *http.Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return httpClient
}

// GetClientByProxy 获取http client
func GetClientByProxy(host string) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://" + host)
		},
	}
	return &http.Client{Transport: tr}
}
