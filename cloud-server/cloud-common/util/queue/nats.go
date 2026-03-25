package queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	jsStreamName    = "SINKING_QUEUE"
	jsSubject       = "SINKING_QUEUE.>"
	jsSubjectPrefix = "SINKING_QUEUE."
)

type NatsQueue struct {
	client        *nats.Conn
	js            nats.JetStreamContext
	unicastSubs   map[string]*nats.Subscription
	broadcastSubs map[string]*nats.Subscription
	mu            sync.RWMutex
}

// ListenReceiveUnicast 监听接收单播消息
// 单播使用 JetStream 消费者组实现，消息持久化存储
func (n *NatsQueue) ListenReceiveUnicast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.unicastSubs == nil {
		n.unicastSubs = make(map[string]*nats.Subscription)
	}
	if oldSub, ok := n.unicastSubs[topic]; ok {
		_ = oldSub.Unsubscribe()
	}
	queueName := "unicast-" + topic
	subject := jsSubjectPrefix + topic
	sub, err := n.js.QueueSubscribe(subject, queueName, func(msg *nats.Msg) {
		if receiveDo(MqMsg{
			RunType: ReceiveUnicast,
			Topic:   msg.Subject,
			Body:    string(msg.Data),
		}) {
			_ = msg.Ack()
		}
	}, nats.Durable(queueName+"-d"))
	if err != nil {
		return err
	}
	n.unicastSubs[topic] = sub
	return nil
}

// ListenReceiveBroadcast 监听接收广播消息
// 广播使用 JetStream 订阅，消息持久化存储
func (n *NatsQueue) ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.broadcastSubs == nil {
		n.broadcastSubs = make(map[string]*nats.Subscription)
	}
	subject := jsSubjectPrefix + topic
	sub, err := n.js.Subscribe(subject, func(msg *nats.Msg) {
		if receiveDo(MqMsg{
			RunType: ReceiveBroadcast,
			Topic:   msg.Subject,
			Body:    string(msg.Data),
		}) {
			_ = msg.Ack()
		}
	})
	if err != nil {
		return err
	}
	n.broadcastSubs[topic] = sub
	return nil
}

// SendUnicast 发送单播消息
// JetStream Publish 实现消息持久化
func (n *NatsQueue) SendUnicast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}
	subject := jsSubjectPrefix + topic
	_, err = n.js.Publish(subject, []byte(body))
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
// JetStream Publish 实现消息持久化
func (n *NatsQueue) SendBroadcast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}
	subject := jsSubjectPrefix + topic
	_, err = n.js.Publish(subject, []byte(body))
	if err != nil {
		return mqMsg, err
	}
	return MqMsg{
		RunType: SendBroadcast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭NATS连接
func (n *NatsQueue) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	// 关闭所有广播订阅
	for _, sub := range n.broadcastSubs {
		_ = sub.Unsubscribe()
	}
	n.broadcastSubs = nil
	// 关闭所有单播订阅
	for _, sub := range n.unicastSubs {
		_ = sub.Unsubscribe()
	}
	n.unicastSubs = nil
	if n.client != nil {
		n.client.Close()
	}
}

// RegisterNatsQueueClient 注册NATS队列实例
func RegisterNatsQueueClient(natsURL string, options ...interface{}) (Queue, error) {
	if natsURL == "" {
		return nil, errors.New("nats url cannot be empty")
	}

	var natsOptions []nats.Option
	for _, opt := range options {
		if natsOpt, ok := opt.(nats.Option); ok {
			natsOptions = append(natsOptions, natsOpt)
		}
	}

	nc, err := nats.Connect(natsURL, natsOptions...)
	if err != nil {
		return nil, errors.New("failed to connect to nats server: " + err.Error())
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, errors.New("failed to create jetstream context: " + err.Error())
	}

	q := &NatsQueue{
		client:        nc,
		js:            js,
		unicastSubs:   make(map[string]*nats.Subscription),
		broadcastSubs: make(map[string]*nats.Subscription),
	}
	// 创建 Stream（幂等：已存在则忽略）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     jsStreamName,
		Subjects: []string{jsSubject},
		MaxBytes: 1024 * 1024 * 1024, // 1GB
		MaxAge:   7 * 24 * time.Hour, // 7天
		Storage:  nats.FileStorage,
	})
	// 等待 Stream 就绪
	_, _ = js.StreamInfo(jsStreamName, nats.Context(ctx))
	return q, nil
}
