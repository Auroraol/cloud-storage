# Pulsar消息队列集成

本模块提供了Apache Pulsar消息队列的集成，用于微服务之间的异步通信。

## 功能特性

- 提供统一的Pulsar客户端管理
- 支持生产者和消费者的创建和管理
- 支持同步和异步消息发送
- 提供消息消费的便捷接口

## 安装和配置

### 1. 安装Pulsar

在开发环境中，可以使用Docker快速启动Pulsar：

```bash
docker run -it -p 6650:6650 -p 8080:8080 --name pulsar apachepulsar/pulsar:3.1.0 bin/pulsar standalone
```

### 2. 添加依赖

确保项目的go.mod中包含Pulsar客户端依赖：

```go
require (
    github.com/apache/pulsar-client-go v0.12.0
)
```

## 使用示例

### 1. 初始化Pulsar管理器

```go
import "github.com/Auroraol/cloud-storage/common/mq"

// 创建Pulsar管理器
manager, err := mq.NewPulsarManager(mq.PulsarConfig{
    URL: "pulsar://localhost:6650",
})
if err != nil {
    log.Fatalf("创建Pulsar管理器失败: %v", err)
}
defer manager.Close()
```

### 2. 发送消息

```go
// 创建消息
message := mq.FileUploadedMessage{
    FileID:      "file-123456",
    FileName:    "example.pdf",
    FileSize:    1024 * 1024 * 5, // 5MB
    ContentType: "application/pdf",
    UserID:      "user-123",
    UploadTime:  time.Now(),
    StoragePath: "/storage/files/user-123/example.pdf",
}

// 发送消息
if err := mq.SendFileUploadedMessage(manager, message); err != nil {
    log.Errorf("发送文件上传消息失败: %v", err)
}
```

### 3. 消费消息

```go
// 启动消费者
err = mq.StartFileUploadedConsumer(manager, func(msg mq.FileUploadedMessage) error {
    log.Infof("收到文件上传消息: %+v", msg)
    // 处理文件上传后的逻辑，如生成缩略图、提取元数据等
    return nil
})
if err != nil {
    log.Errorf("启动文件上传消息消费者失败: %v", err)
}
```

### 4. 直接使用Pulsar客户端API

如果需要更灵活的控制，可以直接使用Pulsar客户端API：

```go
// 创建消费者
consumer, err := manager.CreateConsumer(pulsar.ConsumerOptions{
    Topic:            "my-topic",
    SubscriptionName: "my-subscription",
    Type:             pulsar.Shared,
    Name:             "my-consumer",
})
if err != nil {
    log.Errorf("创建消费者失败: %v", err)
    return
}
defer consumer.Close()

// 接收消息
for {
    msg, err := consumer.Receive(context.Background())
    if err != nil {
        log.Errorf("接收消息失败: %v", err)
        time.Sleep(1 * time.Second)
        continue
    }

    // 处理消息
    log.Infof("收到消息: %s", string(msg.Payload()))
    
    // 确认消息
    consumer.Ack(msg)
}
```

## 在项目中的应用场景

### 1. 文件上传通知

当用户上传文件后，上传服务可以发送消息到Pulsar，通知其他服务进行后续处理，如：
- 存储服务：将文件存储到对象存储
- 缩略图服务：为图片和视频生成缩略图
- 索引服务：更新搜索索引
- 通知服务：通知用户上传完成

### 2. 用户活动日志

记录用户的各种活动，如登录、文件操作、分享等，用于审计和分析。

### 3. 系统事件通知

系统级别的事件通知，如服务启动、停止、配置变更等。

### 4. 异步任务处理

将耗时的任务异步处理，如文件压缩、格式转换、大文件处理等。

## 最佳实践

1. **消息持久化**：使用持久化主题（persistent://）确保消息不会丢失
2. **消息幂等性**：消费者应该处理重复消息的情况
3. **错误处理**：正确处理消息处理失败的情况，使用Nack机制
4. **监控**：监控Pulsar的性能和健康状况
5. **分区主题**：对于高吞吐量的场景，使用分区主题提高并行处理能力 