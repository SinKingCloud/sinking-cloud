package music

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unicode/utf8"
)

type GeCi struct {
	Content   string `json:"content"`
	Info      string `json:"info"`
	Status    int    `json:"status"`
	Charset   string `json:"charset"`
	ErrorCode int    `json:"error_code"`
	Fmt       string `json:"fmt"`
}

func test() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://lyrics.kugou.com/download?ver=1&client=pc&id=93744918&accesskey=B899877B8893EC900F7EDC5F8EB1D247&fmt=krc&charset=utf8", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 Edg/101.0.1210.39")
	req.Header.Set("Cookie", "Hm_lvt_aedee6983d4cfc62f509129360d6bb3d=1650193883,1652176161,1652255250,1652378618; KugooID=413047804; a_id=1014; kg_dfid=4FIE7x44Uh5W1ezhWJ0jFsfR; kg_dfid_collect=d41d8cd98f00b204e9800998ecf8427e; kg_mid=bea091ee0191eb360c522d0bf094d018; kg_mid_temp=bea091ee0191eb360c522d0bf094d018; KuGoo=KugooID=413047804&KugooPwd=F57193B079C921BA59D4D80BDB3B5EE7&NickName=%u6d41%u901d%u4e2d%u6c89%u6ca6&Pic=http://imge.kugou.com/kugouicon/165/20200410/20200410180927579440.jpg&RegState=1&RegFrom=&t=4daac436d7fbf1d470d908bcf9d4af7f6abfff4d25a8d58761df5bec23baad22&t_ts=1652414906&t_key=&a_id=1014&ct=1652414906&UserName=%u006b%u0067%u006f%u0070%u0065%u006e%u0034%u0031%u0033%u0030%u0034%u0037%u0038%u0030%u0034; t=4daac436d7fbf1d470d908bcf9d4af7f6abfff4d25a8d58761df5bec23baad22; UserName=%u006b%u0067%u006f%u0070%u0065%u006e%u0034%u0031%u0033%u0030%u0034%u0037%u0038%u0030%u0034; mid=bea091ee0191eb360c522d0bf094d018; dfid=4FIE7x44Uh5W1ezhWJ0jFsfR; Hm_lpvt_aedee6983d4cfc62f509129360d6bb3d=1652415152")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	body := str.NewStringTool(string(bodyText)).ConvertToString("gbk", "utf-8")
	var t *GeCi
	_ = json.Unmarshal([]byte(body), &t)
	decoded, _ := base64.StdEncoding.DecodeString(t.Content)
	content := decoded[4:]
	key := []int{64, 71, 97, 119, 94, 50, 116, 71, 81, 54, 49, 45, 206, 210, 110, 105}
	var test []byte
	for i := 0; i < len(content); i++ {
		r, _ := utf8.DecodeRune([]byte(string(content[i])))
		test = append(test, byte(int(r)^key[i%16]))
	}
	z, err := zlib.NewReader(bytes.NewReader(test))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(z io.ReadCloser) {
		err := z.Close()
		if err != nil {
			return
		}
	}(z)
	p, err := io.ReadAll(z)
	if err == nil {
		fmt.Println(string(p))
	}
	os.Exit(0)
}
