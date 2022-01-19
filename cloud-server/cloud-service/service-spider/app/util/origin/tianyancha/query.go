package tianyancha

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Search(keyWord string) (*Data, error) {
	url := "https://capi.tianyancha.com/cloud-tempest/app/searchCompany"
	req, _ := http.NewRequest("POST", url, strings.NewReader("{\"sortType\":0,\"pageSize\":20,\"pageNum\":1,\"word\":\""+keyWord+"\",\"allowModifyQuery\":1}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://servicewechat.com/wx9f2867fc22873452/74/page-frame.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	resp, err := GetClient().Do(req)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.State == "ok" {
		return &response.Data, nil
	}
	return nil, errors.New("server error")
}
