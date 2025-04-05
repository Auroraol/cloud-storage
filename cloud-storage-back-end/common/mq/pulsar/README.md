# Pulsar 消息队列客户端

本模块提供了 Apache Pulsar 消息队列的客户端实现，用于微服务之间的异步通信。

## 功能特性

- 提供统一的 Pulsar 客户端管理
- 支持生产者和消费者的创建和管理
- 支持同步和异步消息发送
- 提供消息消费的便捷接口
- 支持 TLS 和认证
- 支持泛型消息处理

## 安装和配置

### 1. 安装 Pulsar

在开发环境中，可以使用 Docker 快速启动 Pulsar：

```bash
docker run -it -p 6650:6650 -p 8080:8080 --name pulsar apachepulsar/pulsar:3.1.0 bin/pulsar standalone
```

### 2. 添加依赖

确保项目的 go.mod 中包含 Pulsar 客户端依赖：

```go
require (
	github.com/apache/pulsar-client-go v0.14.0
)
```

## 使用示例

### 1. 初始化 Pulsar 管理器

```go
import "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/mq/pulsar"

// 创建 Pulsar 管理器
manager, err := pulsar.NewPulsarManager(pulsar.Config{
	URL: "pulsar://localhost:6650",
})
if err != nil {
	log.Fatalf("创建 Pulsar 管理器失败: %v", err)
}
defer manager.Close()
```

### 2. 创建发布者并发送消息

```go
// 创建发布者
publisher, err := pulsar.NewPublisher(manager, pulsar.PubConfig{
	Topic:           "example-topic",
	BatchingEnabled: true,
})
if err != nil {
	log.Fatalf("创建发布者失败: %v", err)
}
defer publisher.Close()

// 创建消息
message := MyMessage{
	ID:   "msg-123",
	Data: "Hello, Pulsar!",
}

// 发送消息
ctx := context.Background()
messageID, err := publisher.SendObject(ctx, message, nil)
if err != nil {
	log.Fatalf("发送消息失败: %v", err)
}

log.Printf("消息已发送，ID: %v", messageID)
```

### 3. 创建订阅者并接收消息

```go
// 创建订阅者
subscriber, err := pulsar.NewSubscriber(manager, pulsar.SubConfig{
	Topic:                      "example-topic",
	SubscriptionName:           "example-subscription",
	SubscriptionType:           "Shared",
	SubscriptionInitialPosition: "Earliest",
	AutoAck:                    false,
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
	// 使用泛型函数反序列化消息
	message, err := pulsar.UnmarshalMessage[MyMessage](msg)
	if err != nil {
		log.Printf("反序列化消息失败: %v", err)
		return err
	}

	// 处理消息
	log.Printf("收到消息: ID=%s, Data=%s", message.ID, message.Data)

	// 返回 nil 表示处理成功
	return nil
})
```

### 4. 异步订阅消息

```go
// 创建上下文（可取消）
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 异步订阅消息
subscriber.SubscribeAsync(ctx, func(msg pulsar.Message) error {
	// 处理消息
	// ...
	return nil
})

// 主程序继续执行其他任务
log.Println("异步订阅已启动，主程序继续执行...")
```

## 配置选项

### Pulsar 管理器配置

```go
pulsar.Config{
	// Pulsar 服务地址，例如 "pulsar://localhost:6650"
	URL: "pulsar://localhost:6650",
	// 操作超时时间
	OperationTimeout: 30 * time.Second,
	// 连接超时时间
	ConnectionTimeout: 30 * time.Second,
	// 是否启用 TLS
	EnableTLS: false,
	// TLS 证书路径
	TLSTrustCertsFilePath: "/path/to/ca.pem",
	// 认证类型和参数
	AuthType: "token",
	AuthParams: map[string]string{
		"token": "your-auth-token",
	},
}
```

### 生产者配置

```go
pulsar.PubConfig{
	// 主题名称
	Topic: "my-topic",
	// 是否启用批量发送
	BatchingEnabled: true,
	// 批量发送最大消息数
	BatchingMaxMessages: 1000,
	// 批量发送最大发布延迟（毫秒）
	BatchingMaxPublishDelay: 10,
	// 压缩类型：None, LZ4, ZLIB, ZSTD
	CompressionType: "LZ4",
	// 发送超时（秒）
	SendTimeout: 30,
	// 最大等待发送的消息数
	MaxPendingMessages: 1000,
}
```

### 消费者配置

```go
pulsar.SubConfig{
	// 主题名称
	Topic: "my-topic",
	// 多主题订阅
	Topics: []string{"topic1", "topic2"},
	// 订阅名称
	SubscriptionName: "my-subscription",
	// 订阅类型：Exclusive, Shared, Failover, KeyShared
	SubscriptionType: "Shared",
	// 消费者名称
	Name: "my-consumer",
	// 初始位置：Latest, Earliest
	SubscriptionInitialPosition: "Earliest",
	// 是否自动确认消息
	AutoAck: false,
	// 未确认消息重新投递延迟（秒）
	NackRedeliveryDelay: 60,
	// 接收队列大小
	ReceiverQueueSize: 1000,
}
```

## 最佳实践

1. **消息持久化**：使用持久化主题（persistent://）确保消息不会丢失
2. **消息幂等性**：消费者应该处理重复消息的情况
3. **错误处理**：正确处理消息处理失败的情况，使用 Nack 机制
4. **资源释放**：使用 defer 确保资源正确释放
5. **上下文管理**：使用 context 管理消息处理的生命周期
6. **批量处理**：对于高吞吐量场景，启用批量发送提高性能
7. **监控**：监控 Pulsar 的性能和健康状况 