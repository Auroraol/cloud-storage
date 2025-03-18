package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/mq/pulsar"
	"github.com/Auroraol/cloud-storage/file_processor/processor"
	pulsarClient "github.com/apache/pulsar-client-go/pulsar"
	"github.com/zeromicro/go-zero/core/conf"
	"go.uber.org/zap"
)

// Config 配置结构
type Config struct {
	LogConfig logx.LogConfig // 日志配置
	SubConfig pulsar.SubConfig
	PubConfig pulsar.PubConfig // 添加发布者配置
}

func main() {
	var configFile = flag.String("f", "file_processor/etc/config.yaml", "the config file")
	flag.Parse()

	// 加载配置
	var c Config
	conf.MustLoad(*configFile, &c)

	// 初始化日志
	c.LogConfig.CustomLevels = map[string]string{
		"processor": "info", // 自定义业务日志级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		panic(err)
	}

	// 初始化 Pulsar 管理器
	if !c.SubConfig.Enabled {
		zap.S().Fatal("Pulsar 未启用，文件处理服务无法工作")
	}

	pulsarManager, err := pulsar.NewPulsarManager(pulsar.Config{
		URL: c.SubConfig.URL,
	})
	if err != nil {
		zap.S().Fatalf("Pulsar 管理器初始化失败: %s", err)
	}
	defer pulsarManager.Close()

	// 创建订阅者
	subscriber, err := pulsar.NewSubscriber(pulsarManager, c.SubConfig)
	if err != nil {
		zap.S().Fatalf("Pulsar 订阅者初始化失败: %s", err)
	}
	defer subscriber.Close()

	// 创建上下文（可取消）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 订阅消息
	go func() {
		err := subscriber.Subscribe(ctx, func(msg pulsarClient.Message) error {
			// 使用泛型函数反序列化消息
			fileMsg, err := pulsar.UnmarshalMessage[pulsar.FileUploadedMessage](msg)
			if err != nil {
				zap.S().Errorf("反序列化消息失败: %v", err)
				return err
			}

			// 处理文件上传消息
			zap.S().Infof("收到文件上传消息: ID=%s, 文件名=%s, 大小=%d, 用户ID=%s, 路径=%s",
				fileMsg.FileID, fileMsg.FileName, fileMsg.FileSize, fileMsg.UserID, fileMsg.StoragePath)

			// 使用处理器处理文件
			err = processor.ProcessFile(fileMsg)
			if err != nil {
				zap.S().Errorf("处理文件失败: %v", err)
				return err
			}

			zap.S().Infof("文件处理完成: %s", fileMsg.FileID)
			return nil
		})

		if err != nil {
			zap.S().Errorf("订阅消息失败: %v", err)
		}
	}()

	zap.S().Info("文件处理服务已启动，等待消息...")

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	zap.S().Info("接收到中断信号，正在关闭服务...")
	cancel() // 取消上下文，停止订阅
	time.Sleep(1 * time.Second)
	zap.S().Info("服务已关闭")
}

// 判断文件类型的辅助函数
func isImageFile(fileName string) bool {
	// 实现判断图片文件的逻辑
	return false
}

func isVideoFile(fileName string) bool {
	// 实现判断视频文件的逻辑
	return false
}

func isDocumentFile(fileName string) bool {
	// 实现判断文档文件的逻辑
	return false
}

// 处理不同类型文件的函数
func processImageFile(msg pulsar.FileUploadedMessage) error {
	// 实现图片处理逻辑
	return nil
}

func processVideoFile(msg pulsar.FileUploadedMessage) error {
	// 实现视频处理逻辑
	return nil
}

func processDocumentFile(msg pulsar.FileUploadedMessage) error {
	// 实现文档处理逻辑
	return nil
}

func processOtherFile(msg pulsar.FileUploadedMessage) error {
	// 实现其他文件处理逻辑
	return nil
}
