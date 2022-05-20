package ocr

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/rpc"
	"io"
	"io/ioutil"
	"net/http"
)

type ControllerOcr struct {
}

// Test 测试
func (ControllerOcr) Test(c *sinking_server.Context) {
	getUrl := func() string {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://api.clwl.online/ocr.php?url=http://captcha.qq.com/getimage?aid=2000201&uin=0&0.9751382343431614", nil)
		if err != nil {
			return ""
		}
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36 Edg/101.0.1210.47")
		resp, err := client.Do(req)
		if err != nil {
			return ""
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		return string(bodyText)
	}
	url := getUrl()
	if url == "" {
		c.Error("fail")
	} else {
		cookie := rpc.Spider.OcrCaptcha(url)
		if cookie == "" {
			c.Error("fail")
		} else {
			c.SuccessWithData("success", cookie)
		}
	}
}
