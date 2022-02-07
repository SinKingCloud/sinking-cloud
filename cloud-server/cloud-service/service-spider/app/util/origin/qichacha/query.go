package qichacha

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Search(token string, keyWord string) (*Result, int) {
	t := strconv.Itoa(int(time.Now().Unix())) + "000"
	url := "https://xcx.qcc.com/mp-weixin/forwardApp/v3/base/advancedSearch?token=" + token + "&t=" + t + "000&pageIndex=1&needGroup=yes&insuredCntStart=&insuredCntEnd=&registCapiBegin=&registCapiEnd=&countyCode=&province=&sortField=&isSortAsc=&searchKey=" + url.QueryEscape(keyWord) + "&searchIndex=default&industryV3="
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://servicewechat.com/wx395200814fcd7599/177/page-frame.html")
	req.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("Qcc-Platform", "mp-weixin")
	req.Header.Set("Qcc-Timestamp", t)
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
	if err != nil || response.Status == 40103 {
		return nil, origin.NotAllow
	}
	if response.Status == 40101 {
		return nil, origin.TokenError
	}
	if response.Status != 200 || response.Result.Status != 200 {
		return nil, origin.Fail
	}
	return &response.Result, origin.Success
}
