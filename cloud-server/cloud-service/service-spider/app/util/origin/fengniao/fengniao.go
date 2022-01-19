package fengniao

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"io"
	"net/http"
	"strings"
)

const URL = "https://www.riskbird.com"

//请求client
var birdClient *http.Client

// GetClient 获取http client
func GetClient() *http.Client {
	if birdClient == nil {
		birdClient = &http.Client{}
	}
	return birdClient
}

// GetJSessionId 生成一个新的session id
func GetJSessionId() string {
	url := URL + "/localstoreageLogin?localLogin=1&r=%2Fweb%2Fquery%2Fsearch"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.riskbird.com")
	req.Header.Set("sec-ch-ua", "\"Microsoft Edge\";v=\"93\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"93\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	resp, err := GetClient().Do(req)
	if err != nil {
		return ""
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	cookie := strings.ReplaceAll(resp.Header.Get("Set-Cookie"), " ", "")
	if cookie == "" {
		return ""
	}
	return str.NetStringTool(cookie).GetBetween("JSESSIONID=", ";")
}
