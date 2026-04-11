package socket

import (
	"errors"
	"hash/maphash"

	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

// getConnectionPool 获取链接池
func (s *service) getConnectionPool(group string) *sinking_websocket.Registry {
	index := maphash.String(s.groupSeed, group) % uint64(len(s.groupShards))
	shard := &s.groupShards[index]
	if value, ok := shard.Load(group); ok {
		return value.(*sinking_websocket.Registry)
	}
	registry := sinking_websocket.NewRegistry()
	actual, _ := shard.LoadOrStore(group, registry)
	return actual.(*sinking_websocket.Registry)
}

// setConnection 设置连接
func (s *service) setConnection(group string, id string, conn *sinking_websocket.Connection) error {
	registry := s.getConnectionPool(group)
	if registry == nil {
		return errors.New("连接池不存在")
	}
	registry.Store(id, conn)
	return nil
}

// getConnection 获取连接
func (s *service) getConnection(group string, id string) (*sinking_websocket.Connection, error) {
	registry := s.getConnectionPool(group)
	if registry == nil {
		return nil, errors.New("连接池不存在")
	}
	connection, exists := registry.Load(id)
	if exists {
		return connection, nil
	}
	return nil, errors.New("连接不存在")
}

// deleteConnection 删除连接
func (s *service) deleteConnection(group string, id string) error {
	registry := s.getConnectionPool(group)
	if registry == nil {
		return errors.New("连接池不存在")
	}
	exists := registry.Delete(id)
	if !exists {
		return errors.New("连接不存在")
	}
	return nil
}
