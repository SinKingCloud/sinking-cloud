package socket

import (
	"hash/maphash"
	"runtime"
	socketQueue "service-socket/app/queue/socket"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Service service接口
type Service interface {
	Server(token string) (*sinking_websocket.Server, error)                                                        //socket服务器
	Name() string                                                                                                  //本机名称
	GenConnectionToken(group string, id string) (string, error)                                                    //生成连接token
	GetConnectionIdAddr(group string, id string) (string, error)                                                   //获取连接地址
	SendToLocalId(group string, id string, messageType int, payload []byte) error                                  //给ID发送消息(本机)
	SendToLocalAll(group string, messageType int, payload []byte) (queued int, dropped int, closed int, err error) //给全体发送消息(本机)
	SendToId(group string, id string, isText bool, data []byte)                                                    //给ID发送消息
	SendToAll(group string, isText bool, data []byte)                                                              //给全体发送消息
}

// Session服务
type service struct {
	name        string
	cache       cache.Interface
	queue       *socketQueue.Queue
	groupSeed   maphash.Seed
	groupShards []sync.Map
}

// NewService 实例化service
func NewService(name string, cache cache.Interface, queue *socketQueue.Queue) *service {
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
	ins := &service{
		name:        name,
		cache:       cache,
		queue:       queue,
		groupSeed:   maphash.MakeSeed(),
		groupShards: make([]sync.Map, shardCount),
	}
	if queue != nil {
		queue.SetHandler(&socketQueue.Handler{
			GetConnectionIdAddr: ins.GetConnectionIdAddr,
			SendToLocalId:       ins.SendToLocalId,
			SendToLocalAll:      ins.SendToLocalAll,
		})
	}
	return ins
}
