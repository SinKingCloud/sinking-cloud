package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadRpc(s)
}

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
}
