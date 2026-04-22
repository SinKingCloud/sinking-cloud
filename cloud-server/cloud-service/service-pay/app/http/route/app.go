package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadRpc(s)
	loadApi(s)
	loadUser(s)
	loadAdmin(s)
}
