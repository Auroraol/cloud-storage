# Pulsar在各服务中的应用示例

本文档提供了在云存储项目各个微服务中应用Pulsar消息队列的示例。

## 配置Pulsar

所有服务都需要在配置文件中添加Pulsar配置：

```yaml
# Pulsar配置
Pulsar:
  Enabled: true              # 是否启用Pulsar
  URL: "pulsar://101.37.165.220:6650" # Pulsar服务地址
  ServiceName: "service-name"  # 服务名称
```

## 1. 上传服务 (upload_service)

上传服务在文件上传成功后发送消息到Pulsar，通知其他服务进行后续处理。

```go
// 在文件上传成功后，发送消息到Pulsar
if l.svcCtx.PulsarManager != nil && l.svcCtx.Config.Pulsar.Enabled {
    // 创建文件上传消息
    fileUploadedMsg := mq.FileUploadedMessage{
        FileID:      strconv.FormatInt(int64(identity), 10),
        FileName:    fileHeader.Filename,
        FileSize:    fileHeader.Size,
        ContentType: fileHeader.Header.Get("Content-Type"),
        UserID:      strconv.FormatInt(userId, 10),
        UploadTime:  time.Now(),
        StoragePath: fileUrl,
    }

    // 发送消息
    err := mq.SendFileUploadedMessage(l.svcCtx.PulsarManager, fileUploadedMsg)
    if err != nil {
        // 只记录日志，不影响上传流程
        zap.S().Warnf("发送文件上传消息失败: %v", err)
    } else {
        zap.S().Infof("文件上传消息发送成功，文件ID: %s", fileUploadedMsg.FileID)
    }
}
```

## 2. 日志服务 (log_service)

日志服务可以使用Pulsar作为日志聚合系统，将各个微服务的日志发送到Pulsar，然后由日志服务进行处理和存储。

```go
// 初始化Pulsar日志输出
if conf.EnablePulsar && conf.PulsarURL != "" && conf.PulsarTopic != "" {
    // 创建Pulsar客户端
    client, err := pulsar.NewClient(pulsar.ClientOptions{
        URL: conf.PulsarURL,
    })
    if err != nil {
        return fmt.Errorf("创建Pulsar客户端失败: %v", err)
    }

    // 创建PulsarCore
    pulsarCore, err := NewPulsarCore(
        client,
        conf.PulsarTopic,
        conf.PulsarServiceName,
        zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
            return lvl >= globalLevel
        }),
        encoder,
    )
    if err != nil {
        return fmt.Errorf("创建Pulsar日志核心失败: %v", err)
    }

    // 添加到cores列表
    cores = append(cores, pulsarCore)
}
```

## 3. 用户中心服务 (user_center)

用户中心服务可以使用Pulsar发送用户活动通知，如用户注册、登录、修改信息等。

```go
// 发送用户活动消息
func SendUserActivityMessage(manager *mq.PulsarManager, userId int64, activity string, details map[string]interface{}) error {
    // 创建消息
    activityMsg := map[string]interface{}{
        "user_id":     userId,
        "activity":    activity,
        "details":     details,
        "timestamp":   time.Now(),
        "ip_address":  "x.x.x.x", // 从请求中获取
        "user_agent":  "xxx",     // 从请求中获取
    }

    // 序列化消息
    payload, err := json.Marshal(activityMsg)
    if err != nil {
        return fmt.Errorf("序列化消息失败: %v", err)
    }

    // 设置消息属性
    properties := map[string]string{
        "event_type": "user_activity",
        "user_id":    strconv.FormatInt(userId, 10),
    }

    // 发送消息
    _, err = manager.SendMessage("persistent://public/default/user-activity", payload, properties)
    if err != nil {
        return fmt.Errorf("发送消息失败: %v", err)
    }

    return nil
}
```

## 4. 分享服务 (share_service)

分享服务可以使用Pulsar发送分享通知，如文件分享、分享链接生成等。

```go
// 发送分享通知消息
func SendShareNotificationMessage(manager *mq.PulsarManager, shareId string, userId int64, fileId string) error {
    // 创建消息
    shareMsg := map[string]interface{}{
        "share_id":    shareId,
        "user_id":     userId,
        "file_id":     fileId,
        "timestamp":   time.Now(),
    }

    // 序列化消息
    payload, err := json.Marshal(shareMsg)
    if err != nil {
        return fmt.Errorf("序列化消息失败: %v", err)
    }

    // 设置消息属性
    properties := map[string]string{
        "event_type": "share_created",
        "user_id":    strconv.FormatInt(userId, 10),
    }

    // 发送消息
    _, err = manager.SendMessage("persistent://public/default/share-notification", payload, properties)
    if err != nil {
        return fmt.Errorf("发送消息失败: %v", err)
    }

    return nil
}
```

## 5. 文件处理服务 (file_processor)

文件处理服务作为消费者，处理其他服务发送的消息，如文件上传通知、分享通知等。

```go
// 启动文件上传消息消费者
err = mq.StartFileUploadedConsumer(pulsarManager, func(msg mq.FileUploadedMessage) error {
    zap.S().Infof("收到文件上传消息: %+v", msg)
    
    // 处理文件上传后的逻辑
    // 1. 生成缩略图
    if isImage(msg.ContentType) {
        err := generateThumbnail(msg.StoragePath, msg.FileID)
        if err != nil {
            zap.S().Errorf("生成缩略图失败: %v", err)
        }
    }
    
    // 2. 提取文件元数据
    metadata, err := extractMetadata(msg.StoragePath, msg.ContentType)
    if err != nil {
        zap.S().Errorf("提取文件元数据失败: %v", err)
    } else {
        // 保存元数据到数据库
        saveMetadata(msg.FileID, metadata)
    }
    
    // 3. 更新搜索索引
    updateSearchIndex(msg.FileID, msg.FileName, metadata)
    
    return nil
})
```

## 消息主题设计

在云存储项目中，我们使用以下主题设计：

1. **文件上传通知**: `persistent://public/default/file-uploaded`
2. **用户活动通知**: `persistent://public/default/user-activity`
3. **分享通知**: `persistent://public/default/share-notification`
4. **系统日志**: `persistent://public/default/system-logs`

## 最佳实践

1. **消息持久化**: 使用持久化主题（persistent://）确保消息不会丢失
2. **消息幂等性**: 消费者应该处理重复消息的情况
3. **错误处理**: 正确处理消息处理失败的情况，使用Nack机制
4. **监控**: 监控Pulsar的性能和健康状况
5. **分区主题**: 对于高吞吐量的场景，使用分区主题提高并行处理能力 