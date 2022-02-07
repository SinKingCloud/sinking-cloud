package commerce

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	pool []string
	lock sync.Mutex
)

func GetProxy() string {
	for {
		if len(pool) > 0 {
			lock.Lock()
			break
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
	defer lock.Unlock()
	host := pool[0]
	pool = append(pool[:0], pool[1:]...)
	return host
}

func ReferProxyTimer(url string) {
	go func(req string) {
		for {
			if len(pool) <= 0 {
				lock.Lock()
				pool = getProxy(req)
				lock.Unlock()
			}
			time.Sleep(1 * time.Second)
		}
	}(url)
}

func getProxy(url string) []string {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := origin.GetClient().Do(req)
	if err != nil {
		return []string{}
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
		return []string{}
	}
	temp := string(body)
	if find := strings.Contains(temp, "false"); find {
		return []string{}
	}
	temp = strings.Replace(temp, "\t", "", -1)
	temp = strings.Replace(temp, "\r", "", -1)
	arr := strings.Split(temp, "\n")
	var tempArr []string
	for _, v := range arr {
		if v != "" {
			tempArr = append(tempArr, v)
		}
	}
	return tempArr
}
