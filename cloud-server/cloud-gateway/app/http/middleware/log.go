package middleware

import (
	"bytes"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/constant"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Log 记录日志
func Log() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//websocket直接通过
		if c.Request.Header.Get("Connection") == "Upgrade" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		//普通请求记录日志
		getHeader := func(header http.Header) string {
			temp := ""
			for k, v := range header {
				temp2 := ""
				for _, v1 := range v {
					temp2 += v1 + ";"
				}
				temp += k + ":" + temp2 + "\n"
			}
			return temp
		}
		requestTime := time.Now()                     //请求时间
		method := c.Request.Method                    //请求方式
		header := getHeader(c.Request.Header.Clone()) //请求头
		requestId := c.GetRequestId()                 //请求id
		requestIp := c.ClientIP(false)                //请求ip
		requestUrl := func(request *http.Request) string {
			scheme := "http://"
			if request.TLS != nil {
				scheme = "https://"
			}
			return scheme + request.Host + c.Request.URL.RequestURI()
		}(c.Request) //请求地址
		//获取请求体
		var ioBody io.ReadCloser
		var body []byte
		if !strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
			ioBody = c.Request.Body
			body, err := ioutil.ReadAll(ioBody)
			if err != nil {
				return
			}
			//放回body，防止后面读取不到
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
		//接管write，方便捕获返回数据
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		//执行业务
		c.Next()
		responseHeader := c.Writer.Header().Clone()
		responseBody := ""
		if responseHeader.Get("Content-Type") == "application/json" {
			responseBody = writer.body.String()
		}
		//拼接日志
		log := &model.SystemLogRequest{
			Id:             strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10),
			RequestId:      requestId,
			RequestIp:      requestIp,
			RequestUrl:     requestUrl,
			RequestMethod:  method,
			RequestHeader:  header,
			RequestBody:    string(body),
			RequestTime:    str.DateTime(requestTime),
			ResponseCode:   c.StatusCode,
			ResponseHeader: getHeader(responseHeader),
			ResponseBody:   responseBody,
			ResponseTime:   int(time.Since(requestTime).Microseconds() / 1000),
			GroupName:      c.GetString(constant.ServerGroupName),
			AppName:        c.GetString(constant.ServerName),
			AppAddr:        c.GetString(constant.ServerAddr),
		}
		producer.NewLogRequest().Build(log).Send()
	})
}

// responseBodyWriter 接管下文write
type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

// Write 输出内容
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
