package middleware

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Log 日志
func Log() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		fmt.Println("start")
		c.Next()
		fmt.Println("")
	}
}
