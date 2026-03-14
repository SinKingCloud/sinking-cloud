package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadRpc(s)
}

// loadMiddleware 中间件
func loadMiddleware(s *sinking_web.Engine) {

}

func loadRpc(s *sinking_web.Engine) {
}
