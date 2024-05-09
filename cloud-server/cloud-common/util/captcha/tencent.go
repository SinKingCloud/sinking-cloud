package captcha

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var Tencent *tencent

type tencent struct {
}

// Check 判断验证码是否通过
func (t *tencent) Check(randStr string, ticket string) bool {
	u := "https://cgi.urlsec.qq.com/index.php?m=check&a=gw_check&callback=url_query&url=https%3A%2F%2Fwww.qq.com%2F" + strconv.Itoa(rand.Intn(999999-111111)+111111) + "&ticket=" + url.QueryEscape(ticket) + "&randstr=" + url.QueryEscape(randStr)
	req, e := http.NewRequest(http.MethodGet, u, nil)
	if e != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("Referer", "https://urlsec.qq.com/check.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, e := io.ReadAll(resp.Body)
	if e != nil || body == nil {
		return false
	}
	type Response struct {
		Data   interface{} `json:"data"`
		ReCode int         `json:"reCode"`
	}
	var response *Response
	e = json.Unmarshal([]byte(t.jsonpToJson(string(body))), &response)
	if e != nil {
		return false
	}
	if response.ReCode == -109 {
		return false
	}
	return true
}

// jsonp转json
func (t *tencent) jsonpToJson(jsonp string) string {
	r := regexp.MustCompile(`\((.*?)\)`)
	matches := r.FindStringSubmatch(jsonp)
	if len(matches) > 1 {
		return matches[1]
	}
	return jsonp
}
