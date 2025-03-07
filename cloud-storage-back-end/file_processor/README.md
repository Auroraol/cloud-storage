# 文件处理服务

文件处理服务是一个独立的微服务，用于处理文件上传后的异步任务。它通过Pulsar消息队列接收文件上传通知，并执行相应的处理逻辑。

## 功能特性

- 接收文件上传通知
- 处理文件元数据
- 生成缩略图
- 文件格式转换
- 内容分析
- 更新搜索索引

## 安装和配置

### 1. 安装依赖

确保已安装Go 1.18或更高版本，并且已安装Pulsar。

### 2. 配置环境变量

复制`.env.example`文件为`.env`，并根据实际环境修改配置：

```bash
cp .env.example .env
```

### 3. 构建服务

```bash
go build -o file_processor
```

### 4. 运行服务

```bash
./file_processor
```

或者指定配置文件路径：

```bash
./file_processor -f /path/to/.env
```

## 架构设计

文件处理服务采用消费者模式，从Pulsar消息队列中消费文件上传消息，并执行相应的处理逻辑。

### 消息流程

1. 上传服务接收文件上传请求
2. 上传服务将文件保存到存储系统
3. 上传服务发送文件上传消息到Pulsar
4. 文件处理服务从Pulsar消费消息
5. 文件处理服务执行相应的处理逻辑
6. 处理完成后，更新文件状态

### 消息格式

文件上传消息格式如下：

```json
{
  "file_id": "123456",
  "file_name": "example.pdf",
  "file_size": 1048576,
  "content_type": "application/pdf",
  "user_id": "user-123",
  "upload_time": "2023-01-01T12:00:00Z",
  "storage_path": "/storage/files/user-123/example.pdf"
}
```

## 扩展功能

可以通过添加新的处理逻辑来扩展文件处理服务的功能，例如：

- 文件病毒扫描
- OCR文本识别
- 图像处理
- 视频转码
- 文档索引 