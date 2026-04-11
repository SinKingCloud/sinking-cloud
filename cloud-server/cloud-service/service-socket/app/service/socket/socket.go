package socket

import (
	"hash/maphash"
	"runtime"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Service service接口
type Service interface {
	Server(token string) (*sinking_websocket.Server, error)                                                        //socket服务器
	Name() string                                                                                                  //本机名称
	GenConnectionToken(group string, id string) (string, error)                                                    //生成连接token
	SendToLocalId(group string, id string, messageType int, payload []byte) error                                  //给ID发送消息(本机)
	SendToLocalAll(group string, messageType int, payload []byte) (queued int, dropped int, closed int, err error) //给全体发送消息(本机)
}

// Session服务
type service struct {
	app         *global.Application
	cache       cache.Interface
	groupSeed   maphash.Seed
	groupShards []sync.Map
}

// NewService 实例化service
func NewService(app *global.Application, cache cache.Interface) *service {
	defaultShards := func() int {
		target := runtime.GOMAXPROCS(0) * 4
		if target <= 0 {
			target = runtime.NumCPU() * 8
		}
		if target <= 0 {
			target = 8
		}
		shards := 1
		for shards < target {
			shards <<= 1
		}
		return shards
	}
	shardCount := defaultShards()
	return &service{
		app:         app,
		cache:       cache,
		groupSeed:   maphash.MakeSeed(),
		groupShards: make([]sync.Map, shardCount),
	}
}
