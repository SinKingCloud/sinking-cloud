package hongdunwang

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-spider/app/util/origin"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
}

func GetInfo() (*Response, int) {
	url := "https://api.clwl.online/ip.php"
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := getClient().Do(req)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return nil, origin.Fail
	}
	fmt.Println("请求体", string(body))
	return &Response{}, origin.Success
}
