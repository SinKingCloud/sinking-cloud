package web_uin

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"io"
	"net/http"
	"strings"
)

// GetLoginQrCode 获取设备锁二维码
func GetLoginQrCode(verifyUrl string) string {
	if !strings.Contains(verifyUrl, "uin=") || !strings.Contains(verifyUrl, "&sig=") {
		return ""
	}
	//解析参数
	uin := str.NewStringTool(verifyUrl).GetBetween("uin=", "&sig=")
	sig := str.NewStringTool(verifyUrl).GetBetween("&sig=", "-1")
	//构建url
	url := "https://ti.qq.com/safe/verify/proxy/domain/oidb.tim.qq.com/v3/oidbinterface/oidb_0xc9e_4?uid=" + uin + "&getqrcode=1&sdkappid=39998&actype=2&bkn="
	req, err := http.NewRequest("POST", url, strings.NewReader("{\"str_dev_auth_token\":\""+sig+"\",\"uint32_flag\":1}"))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Cookie", "domainId=338; uin=o"+uin+"; p_uin=o"+uin)
	req.Header.Set("Origin", "https://ti.qq.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("qname-service", "1104449:131072")
	req.Header.Set("qname-space", "Production")
	req.Header.Set("Referer", verifyUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.124 Safari/537.36 Edg/102.0.1245.41")
	resp, err := http.DefaultClient.Do(req)
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
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//解析返回值
	type result struct {
		Uint32GuaranteeStatus int    `json:"uint32_guarantee_status"`
		StrUrl                string `json:"str_url"`
		ActionStatus          string `json:"ActionStatus"`
		ErrorCode             int    `json:"ErrorCode"`
		ErrorInfo             string `json:"ErrorInfo"`
	}
	var res *result
	err = json.Unmarshal(bytes, &res)
	if err != nil || res == nil {
		return ""
	}
	if res.ErrorCode != 0 {
		return ""
	}
	d := "https://ti.qq.com/safe/scanresult?str_url=" + res.StrUrl + "&envfrom=&verify_id=&verify_scene=&uin=" + uin
	return d
}
