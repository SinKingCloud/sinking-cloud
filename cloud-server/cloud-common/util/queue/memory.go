package queue

import (
	"errors"
)

type MemoryQueue struct {
	taskChan chan string
	closed   bool
}

func (r *MemoryQueue) ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) error {
	go func(taskChan chan string, topic string) {
		for value := range taskChan {
			receiveDo(MqMsg{
				RunType: ReceiveMsg,
				Topic:   topic,
				Body:    value,
			})
		}
	}(r.taskChan, topic)
	return nil
}

func (r *MemoryQueue) SendMsg(topic string, body string) (mqMsg MqMsg, err error) {
	if r.taskChan == nil {
		return mqMsg, errors.New("未初始化")
	}
	if r.closed {
		return mqMsg, errors.New("队列已关闭")
	}

	// 使用select避免在关闭的通道上发送消息时阻塞
	select {
	case r.taskChan <- body:
		return MqMsg{
			RunType: SendMsg,
			Topic:   topic,
			Body:    body,
		}, nil
	default:
		return mqMsg, errors.New("队列已满或已关闭")
	}
}

// Close 关闭内存队列
func (r *MemoryQueue) Close() {
	if !r.closed && r.taskChan != nil {
		r.closed = true
		close(r.taskChan)
	}
}

// RegisterMemoryQueueClient 注册queue实例
func RegisterMemoryQueueClient(num int) (Queue, error) {
	return &MemoryQueue{
		taskChan: make(chan string, num),
		closed:   false,
	}, nil
}
