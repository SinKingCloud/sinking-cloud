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
	Server(token string) (*sinking_websocket.Server, error)
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
