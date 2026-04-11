package route

import (
	rpc2 "service-socket/app/http/rpc"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadRpc(s *sinking_web.Engine) {
	rpc.RegisterHandle(s)
	rpc.RegisterFunc("GenConnectionToken", rpc2.Socket.GenConnectionToken)
	rpc.RegisterFunc("SendToId", rpc2.Socket.SendToId)
	rpc.RegisterFunc("SendToAll", rpc2.Socket.SendToAll)
}
