package middleware

import (
	"bytes"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"io/ioutil"
	"log"
	"net/http"
)

// Log 日志
func Log() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		ioBody := c.Request.Body
		body, err := ioutil.ReadAll(ioBody)
		if err != nil {
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		c.Next()
		log.Println("请求内容", string(body), "响应码", c.StatusCode, "响应头:", c.Writer.Header().Clone(), "响应内容", writer.body.String())
	}
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
