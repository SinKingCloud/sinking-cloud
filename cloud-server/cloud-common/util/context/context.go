package context

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Context 接管上下文方便扩展方法
type Context struct {
	*sinking_web.Context
}

// New 实例化上下文对象
func New(c *sinking_web.Context) *Context {
	return &Context{
		Context: c,
	}
}

// HandleFunc 方法转换
func HandleFunc(handler func(c *Context)) func(ctx *sinking_web.Context) {
	return func(c *sinking_web.Context) {
		handler(New(c))
	}
}
