package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
)

// Client 队列客户端
type Client struct {
	topic             Queue
	name              string
	consumer          chan interface{}
	consumerThreadNum int
	consumerFunc      func(param string) error
	producer          chan interface{}
	producerThreadNum int
	taskChan          chan MqMsg // worker pool 任务队列

	// 任务管理
	consumerTask *job.Task
	producerTask *job.Task
}

// Try 错误捕获实现
func (l *Client) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendTask 发送任务（内部使用）
func (l *Client) SendTask(param interface{}) {
	l.producer <- param
}

// SendUnicastTask 发送单播任务
func (l *Client) SendUnicastTask(param interface{}) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendUnicast(l.name, string(marshal))
	}
}

// SendBroadcastTask 发送广播任务
func (l *Client) SendBroadcastTask(param interface{}) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendBroadcast(l.name, string(marshal))
	}
}

// Stop 停止客户端
func (l *Client) Stop() error {
	// 关闭 worker pool
	if l.taskChan != nil {
		close(l.taskChan)
	}

	var err error
	if l.consumerTask != nil {
		if stopErr := l.consumerTask.Stop(); stopErr != nil {
			err = stopErr
		}
	}
	if l.producerTask != nil {
		if stopErr := l.producerTask.Stop(); stopErr != nil {
			err = stopErr
		}
	}

	// 关闭底层队列连接
	if l.topic != nil {
		l.topic.Close()
	}

	return err
}

// GetConsumerStatus 获取消费者任务状态
func (l *Client) GetConsumerStatus() job.TaskStatus {
	if l.consumerTask == nil {
		return job.TaskStopped
	}
	return l.consumerTask.GetStatus()
}

// GetProducerStatus 获取生产者任务状态
func (l *Client) GetProducerStatus() job.TaskStatus {
	if l.producerTask == nil {
		return job.TaskStopped
	}
	return l.producerTask.GetStatus()
}

// GetConsumerStats 获取消费者任务统计信息
func (l *Client) GetConsumerStats() job.TaskStats {
	if l.consumerTask == nil {
		return job.TaskStats{}
	}
	return l.consumerTask.GetStats()
}

// GetProducerStats 获取生产者任务统计信息
func (l *Client) GetProducerStats() job.TaskStats {
	if l.producerTask == nil {
		return job.TaskStats{}
	}
	return l.producerTask.GetStats()
}

// registerUnicast 注册单播消息队列实例
func registerUnicast(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &Client{
		topic:             queue,
		name:              topicName,
		producerThreadNum: producerThreadNum,
		consumerThreadNum: consumerThreadNum,
		consumerFunc:      consumerFunc,
		taskChan:          make(chan MqMsg, consumerThreadNum*2),
	}

	ins.producer = make(chan interface{}, 10000)

	// 启动生产者任务
	ins.producerTask = job.NewTask().
		SetThread(producerThreadNum).
		SetProducer(func(ctx context.Context, channel chan interface{}) error {
			ins.producer = channel
			<-ctx.Done()
			return ctx.Err()
		}).
		SetConsumer(func(ctx context.Context, param interface{}) error {
			ins.Try(func() {
				body, ok := param.(string)
				if ok {
					_, _ = queue.SendUnicast(topicName, body)
				}
			}, nil)
			return nil
		})

	if err := ins.producerTask.Start(); err != nil {
		return nil, err
	}

	// 预创建 worker pool，固定 consumerThreadNum 个 goroutine
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					ins.consumerFunc(mqMsg.Body)
				}, func(e interface{}) {
					// panic，不影响主流程
				})
			}
		}()
	}

	// 监听单播消息，阻塞投递到 worker pool
	err := queue.ListenReceiveUnicast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg // 满了则阻塞
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "单播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// RegisterMemoryUnicast 注册内存单播消息队列实例
func RegisterMemoryUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisUnicast 注册Redis单播消息队列实例
func RegisterRedisUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsUnicast 注册NATS单播消息队列实例
func RegisterNatsUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error, natsURL string, options ...interface{}) (*Client, error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败: " + err.Error())
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// registerBroadcast 注册广播消息队列实例
func registerBroadcast(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &Client{
		topic:             queue,
		name:              topicName,
		producerThreadNum: producerThreadNum,
		consumerThreadNum: consumerThreadNum,
		consumerFunc:      consumerFunc,
		taskChan:          make(chan MqMsg, consumerThreadNum*2),
	}

	ins.producer = make(chan interface{}, 10000)

	// 启动生产者任务
	ins.producerTask = job.NewTask().
		SetThread(producerThreadNum).
		SetProducer(func(ctx context.Context, channel chan interface{}) error {
			ins.producer = channel
			<-ctx.Done()
			return ctx.Err()
		}).
		SetConsumer(func(ctx context.Context, param interface{}) error {
			ins.Try(func() {
				body, ok := param.(string)
				if ok {
					_, _ = queue.SendBroadcast(topicName, body)
				}
			}, nil)
			return nil
		})

	if err := ins.producerTask.Start(); err != nil {
		return nil, err
	}

	// 预创建 worker pool，固定 consumerThreadNum 个 goroutine
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					ins.consumerFunc(mqMsg.Body)
				}, func(e interface{}) {
					// panic，不影响主流程
				})
			}
		}()
	}

	// 监听广播消息，阻塞投递到 worker pool
	err := queue.ListenReceiveBroadcast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg // 满了则阻塞
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "广播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// RegisterMemoryBroadcast 注册内存广播消息队列实例
func RegisterMemoryBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisBroadcast 注册Redis广播消息队列实例
func RegisterRedisBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*Client, error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsBroadcast 注册NATS广播消息队列实例
func RegisterNatsBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error, natsURL string, options ...interface{}) (*Client, error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败: " + err.Error())
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
