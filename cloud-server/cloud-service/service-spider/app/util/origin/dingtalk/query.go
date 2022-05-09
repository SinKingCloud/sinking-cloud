package dingtalk

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Search(keyWord string) (*Data, int) {
	url := "https://holmes.taobao.com/ding/corp/customer/searchWithSummary"
	req, _ := http.NewRequest("POST", url, strings.NewReader("{\"keyword\":\""+keyWord+"\",\"pageNo\":1,\"pageSize\":20}"))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
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
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, origin.NotAllow
	}
	if response.Code != 0 {
		return nil, origin.TokenError
	}
	return &response.Data, origin.Success
}
