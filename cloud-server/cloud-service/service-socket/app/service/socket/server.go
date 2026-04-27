package socket

import (
	"errors"
	"hash/maphash"
	"net/http"
	"service-socket/app/constant"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// Server ws服务器
func (s *service) Server(token string) (*sinking_websocket.Server, error) {
	group, id, err := s.getConnectionTokenInfo(token)
	if err != nil {
		return nil, err
	}
	server := sinking_websocket.NewServer(
		sinking_websocket.WithConnectionIDResolver(func(*http.Request) (string, error) {
			return id, nil
		}),
		sinking_websocket.WithConnectHandler(func(connection *sinking_websocket.Connection) error {
			if err2 := s.setConnection(group, id, connection); err2 != nil {
				return err2
			}
			if !s.setConnectionId(group, id) {
				_ = s.deleteConnection(group, id, connection)
				return errors.New("写入token失败")
			}
			return nil
		}),
		sinking_websocket.WithPongHandler(func(connection *sinking_websocket.Connection, payload string) {
			s.renewConnectionId(group, id)
		}),
		sinking_websocket.WithDisconnectHandler(func(connection *sinking_websocket.Connection, err error) {
			_ = s.deleteConnection(group, id, connection)
		}),
	)
	return server, nil
}

// Name 本地名称
func (s *service) Name() string {
	return s.name
}

// GenConnectionToken 生成授权连接token
func (s *service) GenConnectionToken(group string, id string) (string, error) {
	token := str.NewStringTool(str.GetSnowWorkIns().GetUuid()).Md5(true)
	primary := group + "|" + id
	if !s.cache.SetWithExpire(constant.CacheNameWithSocketToken+token, primary, constant.CacheTimeWithSocketToken) {
		return "", errors.New("获取连接token失败")
	}
	return token, nil
}

// GetConnectionIdAddr 获取ID的连接地址
func (s *service) GetConnectionIdAddr(group string, id string) (string, error) {
	value, exists := s.cache.Get(constant.CacheNameWithSocketAddr + group + "|" + id)
	if !exists || value == "" {
		return "", errors.New("获取地址失败")
	}
	return value, nil
}

// getConnectionTokenInfo 获取token信息
func (s *service) getConnectionTokenInfo(token string) (group string, id string, err error) {
	key := constant.CacheNameWithSocketToken + token
	group, ok := s.cache.Get(key)
	if !ok || group == "" {
		return "", "", errors.New("该token未授权")
	}
	split := strings.Split(group, "|")
	if len(split) < 2 {
		return "", "", errors.New("该token无效")
	}
	group = split[0]
	id = split[1]
	if group == "" || id == "" {
		return "", "", errors.New("获取token数据失败")
	}
	defer s.cache.Delete(key)
	return group, id, nil
}

// setConnectionId 设置连接ID
func (s *service) setConnectionId(group string, id string) bool {
	key := constant.CacheNameWithSocketAddr + group + "|" + id
	if !s.cache.SetWithExpire(key, s.Name(), constant.CacheTimeWithSocketAddr) {
		return false
	}
	return true
}

// renewConnectionId 刷新连接ID有效期
func (s *service) renewConnectionId(group string, id string) bool {
	buckets := constant.CacheTimeWithSocketAddrRenewMinute
	if buckets > 1 {
		now := time.Now()
		bucket := int(maphash.String(s.groupSeed, group+"\x00"+id) % uint64(buckets))
		minuteSlot := int((now.Unix() / 60) % int64(buckets))
		if minuteSlot != bucket {
			return true
		}
	}
	key := constant.CacheNameWithSocketAddr + group + "|" + id
	if s.cache.SetExpire(key, constant.CacheTimeWithSocketAddr) {
		return true
	}
	if s.setConnectionId(group, id) {
		return true
	}
	return false
}
