package sinking_server

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

const (
	SuccessCode    = 200 //请求成功
	ErrorCode      = 500 //请求失败
	TokenErrorCode = 503 //token认证失败
	NotLoginCode   = 403 //未登录
	NotAllowCode   = 401 //无权限
)

func (c *Context) Response(code int32, message interface{}, data interface{}) {
	c.JSON(http.StatusOK, sinking_web.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (c *Context) Return(code int32, message interface{}, data interface{}) {
	c.Response(code, message, data)
}

func (c *Context) Success(message interface{}) {
	c.Response(SuccessCode, message, nil)
}

func (c *Context) SuccessWithData(message interface{}, data interface{}) {
	c.Response(SuccessCode, message, data)
}

func (c *Context) Error(message interface{}) {
	c.Response(ErrorCode, message, nil)
}

func (c *Context) ErrorWithData(message interface{}, data interface{}) {
	c.Response(ErrorCode, message, data)
}

func (c *Context) TokenError(message interface{}, data interface{}) {
	c.Response(TokenErrorCode, message, data)
}

func (c *Context) NotLogin(message interface{}, data interface{}) {
	c.Response(NotLoginCode, message, data)
}

func (c *Context) NotAllow(message interface{}, data interface{}) {
	c.Response(NotAllowCode, message, data)
}
