package queue

import (
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

var (
	ncInit = make(map[string]*nats.Conn)
	jsInit = make(map[string]nats.JetStreamContext)
	initMu sync.Mutex
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
	}, nats.Durable(queueName+"-d"), nats.AckExplicit())
	if err != nil {
		return err
	}
	n.unicastSubs[topic] = sub
	return nil
}

// ListenReceiveBroadcast 监听接收广播消息
// 广播使用 JetStream 订阅，消息持久化存储，每个订阅者都能收到所有消息
func (n *NatsQueue) ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.broadcastSubs == nil {
		n.broadcastSubs = make(map[string]*nats.Subscription)
	}
	if oldSub, ok := n.broadcastSubs[topic]; ok {
		_ = oldSub.Unsubscribe()
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
	}, nats.Durable("broadcast-"+topic+"-d"), nats.AckExplicit())
	if err != nil {
		return err
	}
	n.broadcastSubs[topic] = sub
	return nil
}

// SendUnicast 发送单播消息（异步，高效）
func (n *NatsQueue) SendUnicast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}
	subject := jsSubjectPrefix + topic
	// 异步发送，不等待 ack
	_, err = n.js.PublishAsync(subject, []byte(body))
	if err != nil {
		return mqMsg, err
	}
	return MqMsg{
		RunType: SendUnicast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// SendBroadcast 发送广播消息（异步，高效）
func (n *NatsQueue) SendBroadcast(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}
	subject := jsSubjectPrefix + topic
	_, err = n.js.PublishAsync(subject, []byte(body))
	if err != nil {
		return mqMsg, err
	}
	return MqMsg{
		RunType: SendBroadcast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭当前实例的订阅（连接由单例持有，不关闭）
func (n *NatsQueue) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, sub := range n.broadcastSubs {
		_ = sub.Unsubscribe()
	}
	n.broadcastSubs = nil
	for _, sub := range n.unicastSubs {
		_ = sub.Unsubscribe()
	}
	n.unicastSubs = nil
}

// RegisterNatsQueueClient 注册NATS队列实例（按 URL 单例）
func RegisterNatsQueueClient(natsURL string, options ...interface{}) (Queue, error) {
	if natsURL == "" {
		return nil, errors.New("nats url cannot be empty")
	}
	initMu.Lock()
	if nc, ok := ncInit[natsURL]; ok {
		initMu.Unlock()
		return &NatsQueue{
			client:        nc,
			js:            jsInit[natsURL],
			unicastSubs:   make(map[string]*nats.Subscription),
			broadcastSubs: make(map[string]*nats.Subscription),
		}, nil
	}
	var natsOptions []nats.Option
	natsOptions = append(natsOptions,
		nats.MaxPingsOutstanding(10),
		nats.ReconnectBufSize(32*1024*1024),
	)
	for _, opt := range options {
		if natsOpt, ok := opt.(nats.Option); ok {
			natsOptions = append(natsOptions, natsOpt)
		}
	}
	nc, err := nats.Connect(natsURL, natsOptions...)
	if err != nil {
		initMu.Unlock()
		return nil, errors.New("failed to connect to nats server: " + err.Error())
	}
	ncInit[natsURL] = nc
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(4096))
	if err != nil {
		nc.Close()
		delete(ncInit, natsURL)
		initMu.Unlock()
		return nil, errors.New("failed to create jetstream context: " + err.Error())
	}
	jsInit[natsURL] = js
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     jsStreamName,
		Subjects: []string{jsSubject},
		MaxBytes: 1024 * 1024 * 1024,
		MaxAge:   7 * 24 * time.Hour,
		Storage:  nats.FileStorage,
	})
	_, _ = js.StreamInfo(jsStreamName)
	initMu.Unlock()
	return &NatsQueue{
		client:        nc,
		js:            js,
		unicastSubs:   make(map[string]*nats.Subscription),
		broadcastSubs: make(map[string]*nats.Subscription),
	}, nil
}
