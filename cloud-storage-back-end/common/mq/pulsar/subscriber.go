package pulsar

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// MessageHandler 消息处理函数类型
type MessageHandler func(msg pulsar.Message) error

// Subscriber Pulsar 消息订阅者
type Subscriber struct {
	manager  *PulsarManager
	consumer pulsar.Consumer
	config   SubConfig
}

// NewSubscriber 创建一个新的消息订阅者
func NewSubscriber(manager *PulsarManager, config SubConfig) (*Subscriber, error) {
	if config.Topic == "" && len(config.Topics) == 0 {
		return nil, errors.New("topic or topics is required")
	}

	if config.SubscriptionName == "" {
		return nil, errors.New("subscription name is required")
	}

	// 创建消费者选项
	options := pulsar.ConsumerOptions{
		SubscriptionName: config.SubscriptionName,
	}

	// 设置主题
	if config.Topic != "" {
		options.Topic = config.Topic
	} else if len(config.Topics) > 0 {
		options.Topics = config.Topics
	}

	// 设置消费者名称
	if config.Name != "" {
		options.Name = config.Name
	}

	// 设置订阅类型
	switch config.SubscriptionType {
	case "Shared":
		options.Type = pulsar.Shared
	case "Failover":
		options.Type = pulsar.Failover
	case "KeyShared":
		options.Type = pulsar.KeyShared
	default:
		options.Type = pulsar.Exclusive
	}

	// 设置初始位置
	switch config.SubscriptionInitialPosition {
	case "Earliest":
		options.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
	default:
		options.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
	}

	// 设置接收队列大小
	if config.ReceiverQueueSize > 0 {
		options.ReceiverQueueSize = config.ReceiverQueueSize
	}

	// 设置未确认消息重新投递延迟
	if config.NackRedeliveryDelay > 0 {
		options.NackRedeliveryDelay = time.Duration(config.NackRedeliveryDelay) * time.Second
	}

	// 创建消费者
	consumer, err := manager.CreateConsumer(options)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		manager:  manager,
		consumer: consumer,
		config:   config,
	}, nil
}

// Close 关闭消费者
func (s *Subscriber) Close() {
	if s.consumer != nil {
		s.consumer.Close()
	}
}

// Receive 接收单条消息
func (s *Subscriber) Receive(ctx context.Context) (pulsar.Message, error) {
	return s.consumer.Receive(ctx)
}

// Ack 确认消息
func (s *Subscriber) Ack(msg pulsar.Message) {
	s.consumer.Ack(msg)
}

// Nack 拒绝消息（将重新投递）
func (s *Subscriber) Nack(msg pulsar.Message) {
	s.consumer.Nack(msg)
}

// Subscribe 订阅消息并处理
func (s *Subscriber) Subscribe(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := s.consumer.Receive(ctx)
			if err != nil {
				// 如果上下文已取消，则退出
				if ctx.Err() != nil {
					return ctx.Err()
				}
				// 其他错误，等待一段时间后重试
				time.Sleep(1 * time.Second)
				continue
			}

			// 处理消息
			err = handler(msg)
			if err != nil {
				// 处理失败，拒绝消息
				s.consumer.Nack(msg)
			} else if !s.config.AutoAck {
				// 处理成功，如果不是自动确认，则手动确认
				s.consumer.Ack(msg)
			}
		}
	}
}

// SubscribeAsync 异步订阅消息并处理
func (s *Subscriber) SubscribeAsync(ctx context.Context, handler MessageHandler) {
	go func() {
		_ = s.Subscribe(ctx, handler)
	}()
}

// UnmarshalMessage 将消息反序列化为对象
func UnmarshalMessage[T any](msg pulsar.Message) (T, error) {
	var obj T
	err := json.Unmarshal(msg.Payload(), &obj)
	if err != nil {
		return obj, errors.Wrap(err, "反序列化消息失败")
	}
	return obj, nil
}

// GetConsumer 获取原始消费者
func (s *Subscriber) GetConsumer() pulsar.Consumer {
	return s.consumer
}
