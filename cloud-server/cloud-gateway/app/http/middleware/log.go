package middleware

import (
	"bytes"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"io/ioutil"
	"net/http"
	"time"
)

// Log 记录日志
func Log() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//websocket直接通过
		if c.Request.Header.Get("Connection") == "Upgrade" {
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
		requestIp := c.GetRequestIp()                 //请求ip
		requestUrl := func(request *http.Request) string {
			scheme := "http://"
			if request.TLS != nil {
				scheme = "https://"
			}
			return scheme + request.Host + c.Request.URL.RequestURI()
		}(c.Request) //请求地址
		//获取请求体
		ioBody := c.Request.Body
		body, err := ioutil.ReadAll(ioBody)
		if err != nil {
			return
		}
		//放回body，防止后面读取不到
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//接管write，方便捕获返回数据
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		//执行业务
		c.Next()
		//拼接日志
		log := &model.SystemLogRequest{
			Id:             str.GetSnowWorkIns().GetId(),
			RequestId:      requestId,
			RequestIp:      requestIp,
			RequestUrl:     requestUrl,
			RequestMethod:  method,
			RequestHeader:  header,
			RequestBody:    string(body),
			RequestTime:    str.DateTime(requestTime),
			ResponseCode:   c.StatusCode,
			ResponseHeader: getHeader(c.Writer.Header().Clone()),
			ResponseBody:   writer.body.String(),
			ResponseTime:   int(time.Since(requestTime).Microseconds() / 1000),
			GroupName:      util.Conf.GetString(config.ConfApplicationGroup),
			AppName:        util.Conf.GetString(config.ConfApplicationName),
			AppAddr:        util.Conf.GetString(config.ConfApplicationAddr) + ":" + util.Conf.GetString(config.ConfApplicationPort),
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
