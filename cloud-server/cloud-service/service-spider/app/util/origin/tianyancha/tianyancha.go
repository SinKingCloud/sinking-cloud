package tianyancha

import (
	"net/http"
)

//请求client
var tycClient *http.Client

// GetClient 获取http client
func GetClient() *http.Client {
	if tycClient == nil {
		tycClient = &http.Client{}
	}
	return tycClient
}
