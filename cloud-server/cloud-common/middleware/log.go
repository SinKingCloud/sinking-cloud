package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

func Log(c *context.Context) {
	if strings.ToLower(c.Request.Header.Get("Connection")) == "upgrade" && strings.ToLower(c.Request.Header.Get("Upgrade")) == "websocket" {
		c.Next()
		return
	}
	var reqBody []byte
	if !strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
		temp, err := io.ReadAll(c.Request.Body)
		if err == nil {
			reqBody = temp
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
	}
	writer := &responseBodyWriter{c.Writer, &bytes.Buffer{}, 0}
	c.Writer = writer
	c.Next()
	// 构建请求信息
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	requestHeader := getHeaderString(c.Request.Header)
	requestBody := string(reqBody)
	// 构建响应信息
	statusLine := fmt.Sprintf("HTTP/%d.%d %d %s", c.Request.ProtoMajor, c.Request.ProtoMinor, writer.statusCode, http.StatusText(writer.statusCode))
	responseHeader := getHeaderString(c.Writer.Header())
	responseBody := writer.body.String()
	// 拼接完整日志（请求行+请求头+请求体 + 空行 + 响应行+响应头+响应体）
	fullLog := fmt.Sprintf("%s\r\n%s\r\n%s\r\n\r\n%s\r\n%s\r\n%s",
		requestLine, requestHeader, requestBody,
		statusLine, responseHeader, responseBody)
	global.App.Log.Println(fullLog)
}

type responseBodyWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func (r *responseBodyWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func getHeaderString(header http.Header) string {
	var sb strings.Builder
	for k, v := range header {
		sb.WriteString(k)
		sb.WriteByte(':')
		sb.WriteString(strings.Join(v, ";"))
		sb.WriteByte('\n')
	}
	return sb.String()
}
