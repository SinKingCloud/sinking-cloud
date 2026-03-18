package queue

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsQueue struct {
	client        *nats.Conn
	unicastSubs   map[string]*nats.Subscription
	broadcastSubs map[string]*nats.Subscription
	mu            sync.RWMutex
}

// ListenReceiveUnicast 监听接收单播消息
// 单播使用 NATS 队列组实现，同一队列名只有一个消费者能收到消息
func (n *NatsQueue) ListenReceiveUnicast(topic string, receiveDo func(mqMsg MqMsg)) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}

	n.mu.Lock()
	if n.unicastSubs == nil {
		n.unicastSubs = make(map[string]*nats.Subscription)
	}

	// 如果已有该主题订阅，先取消
	if oldSub, ok := n.unicastSubs[topic]; ok {
		_ = oldSub.Unsubscribe()
	}

	// 使用队列订阅实现单播（负载均衡）
	// 同一队列名的消费者只会收到其中一条消息
	queueName := "unicast-" + topic
	sub, err := n.client.QueueSubscribe(topic, queueName, func(msg *nats.Msg) {
		receiveDo(MqMsg{
			RunType: ReceiveUnicast,
			Topic:   msg.Subject,
			Body:    string(msg.Data),
		})
	})

	if err != nil {
		n.mu.Unlock()
		return err
	}

	n.unicastSubs[topic] = sub
	n.mu.Unlock()
	return nil
}

// ListenReceiveBroadcast 监听接收广播消息
// 广播使用普通订阅，所有订阅者都能收到消息
func (n *NatsQueue) ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg)) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.broadcastSubs == nil {
		n.broadcastSubs = make(map[string]*nats.Subscription)
	}

	// 使用普通订阅实现广播（所有订阅者都能收到）
	// 每个监听器创建独立的订阅
	sub, err := n.client.Subscribe(topic, func(msg *nats.Msg) {
		receiveDo(MqMsg{
			RunType: ReceiveBroadcast,
			Topic:   msg.Subject,
			Body:    string(msg.Data),
		})
	})

	if err != nil {
		return err
	}

	n.broadcastSubs[topic] = sub
	return nil
}

// SendUnicast 发送单播消息
// NATS 单播使用普通发布，接收方通过队列订阅实现点对点
func (n *NatsQueue) SendUnicast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}

	err = n.client.Publish(topic, []byte(body))
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
// NATS 广播使用普通发布，所有订阅者都能收到
func (n *NatsQueue) SendBroadcast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}

	err = n.client.Publish(topic, []byte(body))
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

	// 转换interface{}为nats.Option
	var natsOptions []nats.Option
	for _, opt := range options {
		if natsOpt, ok := opt.(nats.Option); ok {
			natsOptions = append(natsOptions, natsOpt)
		}
	}

	// 连接到NATS服务器
	nc, err := nats.Connect(natsURL, natsOptions...)
	if err != nil {
		return nil, errors.New("failed to connect to nats server: " + err.Error())
	}

	return &NatsQueue{
		client:      nc,
		unicastSubs: make(map[string]*nats.Subscription),
	}, nil
}
