package qiyoucha

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Search(token string, keyWord string) (*RespBody, int) {
	url := "https://companyinfo.urlsec.qq.com/enterprise"
	req, _ := http.NewRequest("POST", url, strings.NewReader("{\"action\":\"GetQXBCompany\",\"app_type\":0,\"search_word\":\""+keyWord+"\",\"page\":0,\"highlight\":1,\"appid\":\"urlsectestfunction\",\"skey\":\""+token+"\"}"))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Referer", "https://servicewechat.com/wxae1a9aa6bb57b0c7/14/page-frame.html")
	req.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("Qcc-Platform", "mp-weixin")
	req.Header.Set("Qcc-Version", "1.0.0")
	resp, err := origin.GetClient().Do(req)
	if err != nil {
		return nil, origin.Fail
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, origin.Fail
	}
	temp := string(body)
	temp = strings.ReplaceAll(temp, "<em>", "")
	temp = strings.ReplaceAll(temp, "</em>", "")
	var response Response
	err = json.Unmarshal([]byte(temp), &response)
	if err != nil {
		return nil, origin.NotAllow
	}
	if response.RespCode == -4 {
		return nil, origin.TokenError
	}
	if response.RespCode != 0 {
		return nil, origin.TokenError
	}
	return &response.RespBody, origin.Success
}
