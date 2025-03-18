package pulsar

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// ExampleMessage 示例消息结构
type ExampleMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ExamplePublisher 示例：创建发布者并发送消息
func ExamplePublisher() {
	// 创建 Pulsar 管理器
	manager, err := NewPulsarManager(Config{
		URL: "pulsar://101.37.165.220:6650",
	})
	if err != nil {
		log.Fatalf("创建 Pulsar 管理器失败: %v", err)
	}
	defer manager.Close()

	// 创建发布者
	publisher, err := NewPublisher(manager, PubConfig{
		Topic:           "example-topic",
		BatchingEnabled: true,
	})
	if err != nil {
		log.Fatalf("创建发布者失败: %v", err)
	}
	defer publisher.Close()

	// 创建消息
	message := ExampleMessage{
		ID:        "msg-123",
		Content:   "Hello, Pulsar!",
		Timestamp: time.Now(),
	}

	// 发送消息
	ctx := context.Background()
	messageID, err := publisher.SendObject(ctx, message, map[string]string{
		"source": "example",
		"type":   "greeting",
	})
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	log.Printf("消息已发送，ID: %v", messageID)

	// 异步发送消息
	for i := 0; i < 5; i++ {
		message := ExampleMessage{
			ID:        "msg-async-" + string(rune(i)),
			Content:   "Async message " + string(rune(i)),
			Timestamp: time.Now(),
		}

		err := publisher.SendObjectAsync(ctx, message, nil, func(msgID pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
			if err != nil {
				log.Printf("异步发送消息失败: %v", err)
				return
			}
			log.Printf("异步消息已发送，ID: %v", msgID)
		})
		if err != nil {
			log.Printf("准备异步发送消息失败: %v", err)
		}
	}

	// 等待异步消息发送完成
	time.Sleep(1 * time.Second)
}

// ExampleSubscriber 示例：创建订阅者并接收消息
func ExampleSubscriber() {
	// 创建 Pulsar 管理器
	manager, err := NewPulsarManager(Config{
		URL: "pulsar://101.37.165.220:6650",
	})
	if err != nil {
		log.Fatalf("创建 Pulsar 管理器失败: %v", err)
	}
	defer manager.Close()

	// 创建订阅者
	subscriber, err := NewSubscriber(manager, SubConfig{
		Topic:                       "example-topic",
		SubscriptionName:            "example-subscription",
		SubscriptionType:            "Shared",
		SubscriptionInitialPosition: "Earliest",
		AutoAck:                     false,
	})
	if err != nil {
		log.Fatalf("创建订阅者失败: %v", err)
	}
	defer subscriber.Close()

	// 创建上下文（可取消）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 订阅消息
	err = subscriber.Subscribe(ctx, func(msg pulsar.Message) error {
		// 打印消息元数据
		log.Printf("收到消息: ID=%v, PublishTime=%v, Properties=%v",
			msg.ID(), msg.PublishTime(), msg.Properties())

		// 反序列化消息
		var message ExampleMessage
		if err := json.Unmarshal(msg.Payload(), &message); err != nil {
			log.Printf("反序列化消息失败: %v", err)
			return err
		}

		// 处理消息
		log.Printf("消息内容: ID=%s, Content=%s, Timestamp=%v",
			message.ID, message.Content, message.Timestamp)

		// 返回 nil 表示处理成功
		return nil
	})

	if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
		log.Fatalf("订阅消息失败: %v", err)
	}
}

// ExampleAsyncSubscriber 示例：异步订阅消息
func ExampleAsyncSubscriber() {
	// 创建 Pulsar 管理器
	manager, err := NewPulsarManager(Config{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		log.Fatalf("创建 Pulsar 管理器失败: %v", err)
	}
	defer manager.Close()

	// 创建订阅者
	subscriber, err := NewSubscriber(manager, SubConfig{
		Topic:                       "example-topic",
		SubscriptionName:            "example-subscription",
		SubscriptionType:            "Shared",
		SubscriptionInitialPosition: "Earliest",
		AutoAck:                     true,
	})
	if err != nil {
		log.Fatalf("创建订阅者失败: %v", err)
	}
	defer subscriber.Close()

	// 创建上下文（可取消）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 异步订阅消息
	subscriber.SubscribeAsync(ctx, func(msg pulsar.Message) error {
		// 使用泛型函数反序列化消息
		message, err := UnmarshalMessage[ExampleMessage](msg)
		if err != nil {
			log.Printf("反序列化消息失败: %v", err)
			return err
		}

		// 处理消息
		log.Printf("异步收到消息: ID=%s, Content=%s, Timestamp=%v",
			message.ID, message.Content, message.Timestamp)

		return nil
	})

	// 主程序继续执行其他任务
	log.Println("异步订阅已启动，主程序继续执行...")

	// 模拟主程序运行一段时间
	time.Sleep(30 * time.Second)

	// 取消上下文，停止订阅
	cancel()
	log.Println("异步订阅已停止")
}
