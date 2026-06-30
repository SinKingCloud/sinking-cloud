package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// 全局计数器，用于生成唯一实例ID
var instanceCounter int64

type RedisQueue struct {
	client          redis.UniversalClient
	unicastSubs     map[string]*redis.PubSub
	unicastCtx      map[string]context.Context
	unicastCancel   map[string]context.CancelFunc
	broadcastSubs   map[string]*redis.PubSub
	broadcastCtx    context.Context
	broadcastCancel context.CancelFunc
	mu              sync.RWMutex

	// 实例ID，用于生成唯一消费者名称
	instanceID int64
}

// ListenReceiveUnicast 监听接收单播消息
// 单播使用 Redis Stream 实现，接收者需要指定 consumer group
func (r *RedisQueue) ListenReceiveUnicast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	if r.client == nil {
		return errors.New("the redis client is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否已经监听，避免重复监听
	if _, ok := r.unicastCtx[topic]; ok {
		return nil
	}

	// 为这个 topic 创建独立的 context
	ctx, cancel := context.WithCancel(context.Background())
	r.unicastCtx[topic] = ctx
	r.unicastCancel[topic] = cancel

	// 创建 consumer group（如果不存在）
	streamKey := "stream:" + topic
	groupName := "unicast-group"

	// 创建消费者组，使用唯一名称实现负载均衡
	_ = r.client.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()

	// 为这个监听器创建唯一的消费者名称
	consumerName := fmt.Sprintf("consumer-%d-%s", r.instanceID, topic)

	// 启动一个 goroutine 持续读取 Stream
	go func(ctx context.Context, client redis.UniversalClient, streamKey, groupName, consumerName string) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var streams []redis.XStream

				// 优先接管长时间未ACK的pending消息，再读取当前consumer的pending消息和新消息。
				messages, _, err := client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
					Stream:   streamKey,
					Group:    groupName,
					Consumer: consumerName,
					MinIdle:  time.Minute,
					Start:    "0-0",
					Count:    1,
				}).Result()
				if err == nil && len(messages) > 0 {
					streams = []redis.XStream{{
						Stream:   streamKey,
						Messages: messages,
					}}
				} else {
					streams, err = client.XReadGroup(ctx, &redis.XReadGroupArgs{
						Group:    groupName,
						Consumer: consumerName,
						Streams:  []string{streamKey, "0"},
						Count:    1,
					}).Result()
					if err == redis.Nil || len(streams) == 0 || len(streams[0].Messages) == 0 {
						streams, err = client.XReadGroup(ctx, &redis.XReadGroupArgs{
							Group:    groupName,
							Consumer: consumerName,
							Streams:  []string{streamKey, ">"},
							Count:    1,
							Block:    2 * time.Second,
						}).Result()
					}
				}

				if err != nil {
					if err == redis.Nil {
						continue // 超时，继续等待
					}
					return
				}

				for _, stream := range streams {
					for _, msg := range stream.Messages {
						// 解析消息体
						body, ok := msg.Values["body"].(string)
						if !ok {
							continue // 类型不对，跳过
						}
						var mqMsg MqMsg
						if err := json.Unmarshal([]byte(body), &mqMsg); err == nil {
							mqMsg.RunType = ReceiveUnicast
							// 只有回调返回 true 才 ACK
							if receiveDo(mqMsg) {
								client.XAck(ctx, streamKey, groupName, msg.ID)
							}
						}
					}
				}
			}
		}
	}(ctx, r.client, streamKey, groupName, consumerName)

	return nil
}

// ListenReceiveBroadcast 监听接收广播消息
// 广播使用 Redis Pub/Sub 实现
func (r *RedisQueue) ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	if r.client == nil {
		return errors.New("the redis client is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 为广播创建独立的 context
	if r.broadcastCtx == nil || r.broadcastCancel == nil {
		r.broadcastCtx, r.broadcastCancel = context.WithCancel(context.Background())
	}

	// 使用特殊前缀区分广播topic
	broadcastTopic := "broadcast:" + topic

	// 每个监听器创建独立的 PubSub
	sub := r.client.Subscribe(r.broadcastCtx, broadcastTopic)
	if r.broadcastSubs == nil {
		r.broadcastSubs = make(map[string]*redis.PubSub)
	}
	r.broadcastSubs[topic] = sub

	go func(sub *redis.PubSub, ctx context.Context, topic string) {
		defer func() {
			if sub != nil {
				sub.Close()
			}
		}()

		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				receiveDo(MqMsg{
					RunType: ReceiveBroadcast,
					Topic:   topic,
					Body:    msg.Payload,
				})
			}
		}
	}(sub, r.broadcastCtx, topic)
	return nil
}

// SendUnicast 发送单播消息
// 单播使用 Redis Stream 实现
func (r *RedisQueue) SendUnicast(topic string, body string) (mqMsg MqMsg, err error) {
	if r.client == nil {
		return mqMsg, errors.New("the redis client is nil")
	}

	streamKey := "stream:" + topic
	ctx := context.Background()

	// 构建消息
	msg := MqMsg{
		RunType: SendUnicast,
		Topic:   topic,
		Body:    body,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return mqMsg, err
	}

	// 发送到 Stream
	_, err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"body": string(msgBytes),
		},
	}).Result()

	if err != nil {
		return mqMsg, err
	}

	return MqMsg{
		RunType: SendUnicast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// SendBroadcast 发送广播消息
// 广播使用 Redis Pub/Sub 实现
func (r *RedisQueue) SendBroadcast(topic string, body string) (mqMsg MqMsg, err error) {
	if r.client == nil {
		return mqMsg, errors.New("the redis client is nil")
	}

	ctx := context.Background()
	// 使用特殊前缀区分广播topic
	broadcastTopic := "broadcast:" + topic
	result, err := r.client.Publish(ctx, broadcastTopic, body).Result()
	if err != nil {
		return mqMsg, err
	}
	if result <= 0 {
		return mqMsg, errors.New("no subscribers for broadcast")
	}
	return MqMsg{
		RunType: SendBroadcast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭Redis队列连接
func (r *RedisQueue) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 取消所有单播监听
	for _, cancel := range r.unicastCancel {
		cancel()
	}
	r.unicastCtx = nil
	r.unicastCancel = nil

	// 取消广播监听
	if r.broadcastCancel != nil {
		r.broadcastCancel()
		r.broadcastCancel = nil
	}
	// 关闭所有广播订阅
	for _, sub := range r.broadcastSubs {
		_ = sub.Close()
	}
	r.broadcastSubs = nil

	// 关闭所有单播订阅
	for _, sub := range r.unicastSubs {
		_ = sub.Close()
	}
	r.unicastSubs = nil
}

// RegisterRedisQueueClient 注册queue实例
func RegisterRedisQueueClient(RedisClient redis.UniversalClient) (Queue, error) {
	if RedisClient == nil {
		return nil, errors.New("the redis client is nil")
	}
	broadcastCtx, broadcastCancel := context.WithCancel(context.Background())
	return &RedisQueue{
		client:          RedisClient,
		unicastSubs:     make(map[string]*redis.PubSub),
		unicastCtx:      make(map[string]context.Context),
		unicastCancel:   make(map[string]context.CancelFunc),
		broadcastCtx:    broadcastCtx,
		broadcastCancel: broadcastCancel,
		instanceID:      atomic.AddInt64(&instanceCounter, 1),
	}, nil
}
