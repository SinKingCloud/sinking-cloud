package queue

import (
	"encoding/json"
	"errors"
	"sync/atomic"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
)

type consumeTask struct {
	mqMsg MqMsg
}

// UnicastClient 单播消息队列客户端（泛型版）
type UnicastClient[T any] struct {
	topic               Queue
	name                string
	consumer            atomic.Value
	consumerReady       chan struct{}
	consumerReadyClosed atomic.Bool
	taskChan            chan *consumeTask
}

// Try 错误捕获实现
func (l *UnicastClient[T]) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendUnicastTask 发送单播任务（自动序列化）。
func (l *UnicastClient[T]) SendUnicastTask(param T) error {
	if l == nil || l.topic == nil {
		return errors.New("队列未初始化")
	}
	marshal, e := json.Marshal(param)
	if e != nil {
		return e
	}
	_, e = l.topic.SendUnicast(l.name, string(marshal))
	return e
}

// SetConsumer 设置队列消费处理器，支持队列先初始化、服务后注入消费逻辑。
func (l *UnicastClient[T]) SetConsumer(consumerFunc func(param T) error) {
	if consumerFunc == nil {
		return
	}
	l.consumer.Store(consumerFunc)
	if l.consumerReady != nil && l.consumerReadyClosed.CompareAndSwap(false, true) {
		close(l.consumerReady)
	}
}

// waitConsumer 等待业务消费处理器完成绑定。
func (l *UnicastClient[T]) waitConsumer() bool {
	if l.consumerReady != nil {
		<-l.consumerReady
	}
	return l.consumer.Load() != nil
}

// consume 等待消费处理器绑定，并确认消息已进入本地消费队列。
func (l *UnicastClient[T]) consume(mqMsg MqMsg) (ok bool) {
	if !l.waitConsumer() {
		return false
	}
	task := &consumeTask{
		mqMsg: mqMsg,
	}
	l.Try(func() {
		l.taskChan <- task
		ok = true
	}, func(e interface{}) {
		ok = false
	})
	return ok
}

// Stop 停止客户端
func (l *UnicastClient[T]) Stop() error {
	if l.consumerReady != nil && l.consumerReadyClosed.CompareAndSwap(false, true) {
		close(l.consumerReady)
	}
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// RegisterUnicast 注册单播消息队列实例
func RegisterUnicast[T any](client Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*UnicastClient[T], error) {
	if client == nil {
		return nil, errors.New("队列客户端不能为空")
	}
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &UnicastClient[T]{
		topic:         client,
		name:          topicName,
		consumerReady: make(chan struct{}),
		taskChan:      make(chan *consumeTask, consumerThreadNum*100),
	}
	if consumerFunc != nil {
		ins.SetConsumer(consumerFunc)
	}

	// 预创建 worker pool，消费时自动反序列化
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for task := range ins.taskChan {
				ins.Try(func() {
					var param T
					if err := json.Unmarshal([]byte(task.mqMsg.Body), &param); err == nil {
						if value := ins.consumer.Load(); value != nil {
							_ = value.(func(T) error)(param)
						}
					}
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听单播消息
	err := client.ListenReceiveUnicast(topicName, func(mqMsg MqMsg) bool {
		return ins.consume(mqMsg)
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "单播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// BroadcastClient 广播消息队列客户端（泛型版）
type BroadcastClient[T any] struct {
	topic               Queue
	name                string
	consumer            atomic.Value
	consumerReady       chan struct{}
	consumerReadyClosed atomic.Bool
	taskChan            chan *consumeTask
}

// Try 错误捕获实现
func (l *BroadcastClient[T]) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendBroadcastTask 发送广播任务（自动序列化）。
func (l *BroadcastClient[T]) SendBroadcastTask(param T) error {
	if l == nil || l.topic == nil {
		return errors.New("队列未初始化")
	}
	marshal, e := json.Marshal(param)
	if e != nil {
		return e
	}
	_, e = l.topic.SendBroadcast(l.name, string(marshal))
	return e
}

// SetConsumer 设置队列消费处理器，支持队列先初始化、服务后注入消费逻辑。
func (l *BroadcastClient[T]) SetConsumer(consumerFunc func(param T) error) {
	if consumerFunc == nil {
		return
	}
	l.consumer.Store(consumerFunc)
	if l.consumerReady != nil && l.consumerReadyClosed.CompareAndSwap(false, true) {
		close(l.consumerReady)
	}
}

// waitConsumer 等待业务消费处理器完成绑定。
func (l *BroadcastClient[T]) waitConsumer() bool {
	if l.consumerReady != nil {
		<-l.consumerReady
	}
	return l.consumer.Load() != nil
}

// consume 等待消费处理器绑定，并确认消息已进入本地消费队列。
func (l *BroadcastClient[T]) consume(mqMsg MqMsg) (ok bool) {
	if !l.waitConsumer() {
		return false
	}
	task := &consumeTask{
		mqMsg: mqMsg,
	}
	l.Try(func() {
		l.taskChan <- task
		ok = true
	}, func(e interface{}) {
		ok = false
	})
	return ok
}

// Stop 停止客户端
func (l *BroadcastClient[T]) Stop() error {
	if l.consumerReady != nil && l.consumerReadyClosed.CompareAndSwap(false, true) {
		close(l.consumerReady)
	}
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// RegisterBroadcast 注册广播消息队列实例
func RegisterBroadcast[T any](client Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*BroadcastClient[T], error) {
	if client == nil {
		return nil, errors.New("队列客户端不能为空")
	}
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &BroadcastClient[T]{
		topic:         client,
		name:          topicName,
		consumerReady: make(chan struct{}),
		taskChan:      make(chan *consumeTask, consumerThreadNum*100),
	}
	if consumerFunc != nil {
		ins.SetConsumer(consumerFunc)
	}

	// 预创建 worker pool，消费时自动反序列化
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for task := range ins.taskChan {
				ins.Try(func() {
					var param T
					if err := json.Unmarshal([]byte(task.mqMsg.Body), &param); err == nil {
						if value := ins.consumer.Load(); value != nil {
							_ = value.(func(T) error)(param)
						}
					}
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听广播消息
	err := client.ListenReceiveBroadcast(topicName, func(mqMsg MqMsg) bool {
		return ins.consume(mqMsg)
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "广播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// RegisterMemoryUnicast 注册内存单播消息队列实例
func RegisterMemoryUnicast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*UnicastClient[T], error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return RegisterUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisUnicast 注册Redis单播消息队列实例
func RegisterRedisUnicast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*UnicastClient[T], error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return RegisterUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsUnicast 注册NATS单播消息队列实例
func RegisterNatsUnicast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error, natsURL string, options ...interface{}) (*UnicastClient[T], error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败: " + err.Error())
	}
	return RegisterUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterMemoryBroadcast 注册内存广播消息队列实例
func RegisterMemoryBroadcast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*BroadcastClient[T], error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return RegisterBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisBroadcast 注册Redis广播消息队列实例
func RegisterRedisBroadcast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*BroadcastClient[T], error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return RegisterBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsBroadcast 注册NATS广播消息队列实例
func RegisterNatsBroadcast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error, natsURL string, options ...interface{}) (*BroadcastClient[T], error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败: " + err.Error())
	}
	return RegisterBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
