package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMq struct {
	endPoints   []string
	producerIns rocketmq.Producer
	consumerIns rocketmq.PushConsumer
}

// RegisterRocketProducerMust 注册并启动生产者接口实现
func RegisterRocketProducerMust(endPoints []string, groupName string, retry int) (client MqProducer) {
	var err error
	client, err = RegisterRocketMqProducer(endPoints, groupName, retry)
	if err != nil {
		panic(err)
	}
	return client
}

// RegisterRocketConsumerMust 注册消费者
func RegisterRocketConsumerMust(endPoints []string, groupName string, allListen bool) (client MqConsumer) {
	var err error
	client, err = RegisterRocketMqConsumer(endPoints, groupName, allListen)
	if err != nil {
		panic(err)
	}
	return client
}

// SendMsg 按字符串类型生产数据
func (r *RocketMq) SendMsg(topic string, tag string, body string) (mqMsg MqMsg, err error) {
	return r.SendByteMsg(topic, tag, []byte(body))
}

// SendByteMsg 生产数据
func (r *RocketMq) SendByteMsg(topic string, tag string, body []byte) (mqMsg MqMsg, err error) {
	if r.producerIns == nil {
		return mqMsg, errors.New("RocketMq producer not register")
	}
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	if tag != "" {
		msg.WithTag(tag)
	}
	result, err := r.producerIns.SendSync(context.Background(), msg)

	if err != nil {
		return
	}
	if result.Status != primitive.SendOK {
		return mqMsg, errors.New(fmt.Sprintf("RocketMq producer send msg error status:%v", result.Status))
	}

	mqMsg = MqMsg{
		RunType: SendMsg,
		Topic:   topic,
		MsgId:   result.MsgID,
		Body:    body,
	}
	return mqMsg, nil
}

// ListenReceiveMsgDo 消费数据
func (r *RocketMq) ListenReceiveMsgDo(topic string, tag string, receiveDo func(mqMsg MqMsg)) (err error) {
	if r.consumerIns == nil {
		return errors.New("RocketMq consumer not register")
	}
	if tag == "" {
		tag = "*"
	}
	selector := &consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	}
	err = r.consumerIns.Subscribe(topic, *selector, func(ctx context.Context, megs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, item := range megs {
			receiveDo(MqMsg{
				RunType: ReceiveMsg,
				Topic:   item.Topic,
				MsgId:   item.MsgId,
				Body:    item.Body,
			})
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return err
	}
	err = r.consumerIns.Start()
	if err != nil {
		err := r.consumerIns.Unsubscribe(topic)
		if err != nil {
			return err
		}
		return err
	}
	return
}

// RegisterRocketMqProducer 注册rocketmq生产者
func RegisterRocketMqProducer(endPoints []string, groupName string, retry int) (mqIns *RocketMq, err error) {
	addr, err := primitive.NewNamesrvAddr(endPoints...)
	if err != nil {
		return nil, err
	}
	mqIns = &RocketMq{
		endPoints: endPoints,
	}

	if retry <= 0 {
		retry = 0
	}

	mqIns.producerIns, err = rocketmq.NewProducer(
		producer.WithNameServer(addr),
		producer.WithRetry(retry),
		producer.WithGroupName(groupName),
	)

	if err != nil {
		return nil, err
	}

	err = mqIns.producerIns.Start()
	if err != nil {
		return nil, err
	}

	return mqIns, nil
}

// RegisterRocketMqConsumer 注册rocketmq消费者
func RegisterRocketMqConsumer(endPoints []string, groupName string, allListen bool) (mqIns *RocketMq, err error) {
	addr, err := primitive.NewNamesrvAddr(endPoints...)
	if err != nil {
		return nil, err
	}
	mqIns = &RocketMq{
		endPoints: endPoints,
	}
	consumerModel := consumer.Clustering
	if allListen {
		consumerModel = consumer.BroadCasting
	}
	mqIns.consumerIns, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer(addr),
		consumer.WithConsumerModel(consumerModel),
		consumer.WithGroupName(groupName),
	)

	if err != nil {
		return nil, err
	}

	return mqIns, nil
}
