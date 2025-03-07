package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
)

// 示例消息结构
type FileUploadedMessage struct {
	FileID      string    `json:"file_id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	ContentType string    `json:"content_type"`
	UserID      string    `json:"user_id"`
	UploadTime  time.Time `json:"upload_time"`
	StoragePath string    `json:"storage_path"`
}

// 文件上传主题
const TopicFileUploaded = "persistent://public/default/file-uploaded"

// 发送文件上传消息
func SendFileUploadedMessage(manager *PulsarManager, message FileUploadedMessage) error {
	// 序列化消息
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 设置消息属性
	properties := map[string]string{
		"event_type": "file_uploaded",
		"user_id":    message.UserID,
	}

	// 发送消息
	_, err = manager.SendMessage(TopicFileUploaded, payload, properties)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// 启动文件上传消息消费者
func StartFileUploadedConsumer(manager *PulsarManager, handler func(FileUploadedMessage) error) error {
	// 创建消费者
	consumer, err := manager.CreateConsumer(pulsar.ConsumerOptions{
		Topic:            TopicFileUploaded,
		SubscriptionName: "file-processor-subscription",
		Type:             pulsar.Shared,
		Name:             "file-processor",
	})
	if err != nil {
		return fmt.Errorf("创建消费者失败: %v", err)
	}

	// 启动消费循环
	go func() {
		defer consumer.Close()

		for {
			// 接收消息
			msg, err := consumer.Receive(context.Background())
			if err != nil {
				zap.S().Errorf("接收消息失败: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// 解析消息
			var uploadMsg FileUploadedMessage
			if err := json.Unmarshal(msg.Payload(), &uploadMsg); err != nil {
				zap.S().Errorf("解析消息失败: %v", err)
				consumer.Nack(msg)
				continue
			}

			// 处理消息
			if err := handler(uploadMsg); err != nil {
				zap.S().Errorf("处理消息失败: %v", err)
				consumer.Nack(msg)
			} else {
				consumer.Ack(msg)
			}
		}
	}()

	return nil
}

// 示例：如何在上传服务中使用
func ExampleUsage() {
	// 创建Pulsar管理器
	manager, err := NewPulsarManager(PulsarConfig{
		URL: "pulsar://localhost:6650",
	})
	if err != nil {
		zap.S().Fatalf("创建Pulsar管理器失败: %v", err)
	}
	defer manager.Close()

	// 发送文件上传消息示例
	message := FileUploadedMessage{
		FileID:      "file-123456",
		FileName:    "example.pdf",
		FileSize:    1024 * 1024 * 5, // 5MB
		ContentType: "application/pdf",
		UserID:      "user-123",
		UploadTime:  time.Now(),
		StoragePath: "/storage/files/user-123/example.pdf",
	}

	if err := SendFileUploadedMessage(manager, message); err != nil {
		zap.S().Errorf("发送文件上传消息失败: %v", err)
	}

	// 启动文件上传消息消费者示例
	err = StartFileUploadedConsumer(manager, func(msg FileUploadedMessage) error {
		zap.S().Infof("收到文件上传消息: %+v", msg)
		// 处理文件上传后的逻辑，如生成缩略图、提取元数据等
		return nil
	})
	if err != nil {
		zap.S().Errorf("启动文件上传消息消费者失败: %v", err)
	}

	// 保持程序运行
	select {}
}
