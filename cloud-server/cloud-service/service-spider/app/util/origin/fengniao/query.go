package fengniao

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// SearchCompany 搜索结果
type SearchCompany struct {
	Name             string
	ArtificialPerson string
}

// Search 搜索关键词企业信息
func Search(keyword string) ([]*SearchCompany, int) {
	sessionId := GetJSessionId()
	if sendPostData(keyword, sessionId) {
		return getData(sessionId)
	}
	return nil, origin.Fail
}

// sendPostData 发送关键词信息
func sendPostData(keyWord string, sessionId string) bool {
	url := URL + "/web/query/search"
	req, err := http.NewRequest("POST", url, strings.NewReader("searchQueryType=1&inputQueryType=1&inputKeyword=&inputIdCard=&keyword="+keyWord+"&idCard="))
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "suspend=1; JSESSIONID="+sessionId+";")
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
	resp, err2 := origin.GetClient().Do(req)
	if err2 != nil {
		return false
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	return true
}

// getData 获取搜索结果
func getData(sessionId string) ([]*SearchCompany, int) {
	url := "https://www.riskbird.com/web/query/search"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "suspend=1; JSESSIONID="+sessionId+";")
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
	if strings.Contains(temp, "风鸟查询平台-登录") {
		return nil, origin.TokenError
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(temp))
	if err != nil {
		return nil, origin.Fail
	}
	var company []*SearchCompany
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		onclick, exists := s.Find("a").Attr("onclick")
		if !exists {
			return
		}
		span := s.Find("label span").Text()
		if span == "" {
			return
		}
		onclick = str.NetStringTool(onclick).GetBetween("query('", "','0')")
		span = strings.ReplaceAll(span, "法定代表人:  ", "")
		company = append(company, &SearchCompany{
			Name:             onclick,
			ArtificialPerson: span,
		})
	})
	return company, origin.Success
}
