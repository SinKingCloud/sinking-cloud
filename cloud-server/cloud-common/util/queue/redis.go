package queue

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

type RedisQueue struct {
	client *redis.Pool
}

func (r *RedisQueue) ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) error {
	if r.client == nil {
		return errors.New("the redis client is nil")
	}
	c := redis.PubSubConn{Conn: r.client.Get()}
	e := c.Subscribe(topic)
	if e != nil {
		return e
	}
	go func(c redis.PubSubConn, receiveDo func(mqMsg MqMsg)) {
		for {
			receive := c.Receive()
			if receive == nil {
				continue
			}
			switch res := receive.(type) {
			case redis.Message:
				receiveDo(MqMsg{
					RunType: ReceiveMsg,
					Topic:   res.Channel,
					Body:    string(res.Data),
				})
			case redis.Subscription:
				continue
			case error:
				continue
			}
		}
	}(c, receiveDo)
	return nil
}

func (r *RedisQueue) SendMsg(topic string, body string) (mqMsg MqMsg, err error) {
	if r.client == nil {
		return mqMsg, errors.New("the redis client is nil")
	}
	do, err := redis.Int64(r.client.Get().Do("Publish", topic, body))
	if err != nil {
		return mqMsg, err
	}
	if do <= 0 {
		return mqMsg, errors.New("the topic is not exists")
	}
	return MqMsg{
		RunType: SendMsg,
		Topic:   topic,
		Body:    body,
	}, nil
}

// RegisterRedisQueueClient 注册queue实例
func RegisterRedisQueueClient(RedisClient *redis.Pool) (Queue, error) {
	if RedisClient == nil {
		return nil, errors.New("the redis client is nil")
	}
	return &RedisQueue{client: RedisClient}, nil
}
