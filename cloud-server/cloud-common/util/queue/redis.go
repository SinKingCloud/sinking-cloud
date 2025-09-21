package queue

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client       *redis.Client
	subscription *redis.PubSub
	ctx          context.Context
	cancel       context.CancelFunc
}

func (r *RedisQueue) ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) error {
	if r.client == nil {
		return errors.New("the redis client is nil")
	}

	// 如果已有订阅，先关闭
	if r.subscription != nil {
		_ = r.subscription.Close()
	}
	if r.cancel != nil {
		r.cancel()
	}

	// 创建新的context
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.subscription = r.client.Subscribe(r.ctx, topic)

	go func() {
		defer func() {
			if r.subscription != nil {
				r.subscription.Close()
			}
		}()

		ch := r.subscription.Channel()
		for {
			select {
			case <-r.ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				receiveDo(MqMsg{
					RunType: ReceiveMsg,
					Topic:   msg.Channel,
					Body:    msg.Payload,
				})
			}
		}
	}()
	return nil
}

func (r *RedisQueue) SendMsg(topic string, body string) (mqMsg MqMsg, err error) {
	if r.client == nil {
		return mqMsg, errors.New("the redis client is nil")
	}
	ctx := context.Background()
	result, err := r.client.Publish(ctx, topic, body).Result()
	if err != nil {
		return mqMsg, err
	}
	if result <= 0 {
		return mqMsg, errors.New("the topic is not exists")
	}
	return MqMsg{
		RunType: SendMsg,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭Redis队列连接
func (r *RedisQueue) Close() {
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
	if r.subscription != nil {
		_ = r.subscription.Close()
		r.subscription = nil
	}
	// 注意：这里不关闭client，因为它可能被其他地方使用
}

// RegisterRedisQueueClient 注册queue实例
func RegisterRedisQueueClient(RedisClient *redis.Client) (Queue, error) {
	if RedisClient == nil {
		return nil, errors.New("the redis client is nil")
	}
	return &RedisQueue{client: RedisClient}, nil
}
