package hongdunwang

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

//请求client
var httpClient *http.Client

// getClient 获取http client
func getClient() *http.Client {
	if httpClient == nil {
		proxy := func(r *http.Request) (*url.URL, error) {
			fmt.Println(r.URL.Scheme)
			return url.Parse("http://110.86.179.101:64257")
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxy,
		}
		httpClient = &http.Client{Transport: tr}
	}
	return httpClient
}
