package queue

import (
	"errors"
	"github.com/nats-io/nats.go"
)

type NatsQueue struct {
	client       *nats.Conn
	subscription *nats.Subscription
}

func (n *NatsQueue) ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) error {
	if n.client == nil {
		return errors.New("the nats client is nil")
	}

	// 如果已有订阅，先取消订阅
	if n.subscription != nil {
		_ = n.subscription.Unsubscribe()
	}

	// 订阅主题
	sub, err := n.client.Subscribe(topic, func(msg *nats.Msg) {
		receiveDo(MqMsg{
			RunType: ReceiveMsg,
			Topic:   msg.Subject,
			Body:    string(msg.Data),
		})
	})

	if err != nil {
		return err
	}

	n.subscription = sub
	return nil
}

func (n *NatsQueue) SendMsg(topic string, body string) (mqMsg MqMsg, err error) {
	if n.client == nil {
		return mqMsg, errors.New("the nats client is nil")
	}

	err = n.client.Publish(topic, []byte(body))
	if err != nil {
		return mqMsg, err
	}

	return MqMsg{
		RunType: SendMsg,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭NATS连接
func (n *NatsQueue) Close() {
	if n.subscription != nil {
		_ = n.subscription.Unsubscribe()
		n.subscription = nil
	}
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

	return &NatsQueue{client: nc}, nil
}
