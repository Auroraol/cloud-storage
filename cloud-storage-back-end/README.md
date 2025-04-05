# 云存储系统

## 系统架构

集成后的系统架构由以下几个部分组成：

- **现有项目**：由Nginx服务器和端口9090上的后端服务组成
- **云存储系统**：基于微服务架构设计的云存储平台
  - 用户中心服务：负责用户认证与管理
  - 上传服务：处理文件上传与存储
  - 分享服务：提供文件分享功能
  - 日志服务：记录系统操作日志
- **基础设施服务**：
  - MySQL：数据库服务
  - Redis：缓存服务
  - ETCD：配置中心
- **日志收集**：
  - zap: 日志库

## 部署准备

### 环境要求

- Docker 19.03.0+
- Docker Compose 1.27.0+
- Linux系统（推荐CentOS 7或Ubuntu 18.04+）
- 4GB+ RAM
- 20GB+ 磁盘空间

### 目录结构

```
/root/cloud-storage/           # 项目根目录
├── deploy_integrated.sh       # 集成部署脚本
├── update_services.sh         # 服务更新脚本
├── integrated_docker-compose.yml  # 集成的Docker Compose配置
├── integrated_nginx.conf      # 集成的Nginx配置
├── prometheus.yml             # Prometheus配置
└── .env                       # 环境变量配置
```

### 配置阿里云OSS

在部署前，您需要在`.env`文件中配置阿里云OSS参数：

```
OSS_ACCESS_KEY_ID=您的AccessKeyID
OSS_ACCESS_KEY_SECRET=您的AccessKeySecret
OSS_BUCKET_NAME=您的存储桶名称
OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
```

## 部署步骤

### 初始部署

1. 创建项目目录并复制所有文件到该目录：

```bash
mkdir -p /root/cloud-storage
cd /root/cloud-storage
# 复制所有配置文件到当前目录
```

2. 添加脚本执行权限：

```bash
chmod +x deploy_integrated.sh
chmod +x update_services.sh
```

3. 执行部署脚本：

```bash
./deploy_integrated.sh
```

部署脚本会自动完成以下任务：
- 检查Docker和Docker Compose安装情况
- 检查现有后端服务是否正常运行
- 创建必要的目录结构
- 配置Nginx
- 启动所有服务
- 显示访问地址

### 更新服务

如果需要更新配置或服务，可以使用更新脚本：

```bash
./update_services.sh
```

更新脚本会执行以下操作：
- 备份当前配置
- 更新配置文件
- 重启服务
- 显示服务状态

# 服务启动顺序

## 第一阶段：RPC 服务启动

1. **上传服务 RPC**
2. **日志服务 RPC**
3. **用户中心 RPC**

## 第二阶段：API 服务启动

1. **用户中心 API**
2. **上传服务 API** 
3. **日志服务 API**
4. **分享服务 API**

# 云存储系统 API 文档

## 基础信息

- **基础URL**：http://101.37.165.220/api/
- **认证方式**：JWT认证（除特殊说明的接口外均需要认证）
- **数据格式**：JSON

### 认证头格式

```
Authorization: Bearer {token}
```

## 目录

1. [用户中心服务](#用户中心服务)
2. [上传服务](#上传服务)
3. [分享服务](#分享服务)
4. [日志服务](#日志服务)

## 用户中心服务

用户中心服务负责用户认证与管理，以及用户文件和文件夹的管理。

### 认证相关接口

#### 账号密码登录

- **接口URL**：`/user/oauth/login`
- **请求方式**：POST
- **认证要求**：无需认证
- **请求参数**：

```json
{
  "username": "用户名/邮箱",
  "password": "密码"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "用户ID",
      "username": "用户名",
      "email": "邮箱",
      "avatar": "头像URL"
    }
  }
}
```

#### 手机号登录/注册

- **接口URL**：`/user/oauth/login/mobile`
- **请求方式**：POST
- **认证要求**：无需认证
- **请求参数**：

```json
{
  "mobile": "手机号",
  "code": "验证码"
}
```

- **响应示例**：与账号密码登录相同

#### 账号密码注册

- **接口URL**：`/user/oauth/register`
- **请求方式**：POST
- **认证要求**：无需认证
- **请求参数**：

```json
{
  "username": "用户名",
  "password": "密码",
  "email": "邮箱",
  "code": "验证码"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "注册成功"
}
```

#### 验证码发送

- **接口URL**：`/user/oauth/send`
- **请求方式**：POST
- **认证要求**：无需认证
- **请求参数**：

```json
{
  "email": "邮箱",
  "type": "验证码类型: register/login/reset"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "验证码已发送"
}
```

### 用户信息相关接口

#### 刷新Authorization

- **接口URL**：`/user/refresh/authorization`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "refresh_token": "刷新令牌"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "令牌刷新成功",
  "data": {
    "token": "新的访问令牌",
    "refresh_token": "新的刷新令牌"
  }
}
```

#### 更换头像

- **接口URL**：`/user/user/avatar/update`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "avatar_url": "头像URL"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "头像更新成功"
}
```

#### 获取用户信息

- **接口URL**：`/user/user/detail`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：无
- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": "用户ID",
    "username": "用户名",
    "email": "邮箱",
    "mobile": "手机号",
    "avatar": "头像URL",
    "storage_used": "已使用存储空间(字节)",
    "storage_max": "最大存储空间(字节)"
  }
}
```

#### 修改用户信息

- **接口URL**：`/user/user/info/update`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "username": "新用户名",
  "email": "新邮箱"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "用户信息更新成功"
}
```

#### 修改密码

- **接口URL**：`/user/user/password/update`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "密码修改成功"
}
```

### 文件与文件夹管理接口

#### 用户文件删除

- **接口URL**：`/user/user/file/delete`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "file_id": "文件ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件删除成功"
}
```

#### 用户文件和文件夹列表

- **接口URL**：`/user/user/file/folder/list`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "parent_id": "父文件夹ID，根目录为0",
  "page": 1,
  "page_size": 20
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total": 10,
    "list": [
      {
        "id": "ID",
        "name": "文件名/文件夹名",
        "ext": "文件扩展名",
        "size": "文件大小(字节)",
        "type": "类型(file/folder)",
        "path": "文件路径",
        "parent_id": "父文件夹ID",
        "created_at": "创建时间",
        "updated_at": "更新时间"
      }
    ]
  }
}
```

#### 用户文件列表

- **接口URL**：`/user/user/file/list`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "page": 1,
  "page_size": 20
}
```

- **响应示例**：与文件和文件夹列表相同

#### 用户文件移动

- **接口URL**：`/user/user/file/move`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "file_id": "文件ID",
  "parent_id": "目标文件夹ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件移动成功"
}
```

#### 用户文件名称修改

- **接口URL**：`/user/user/file/name/update`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "file_id": "文件ID",
  "name": "新文件名"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件名修改成功"
}
```

#### 搜索用户文件和文件夹

- **接口URL**：`/user/user/file/search`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "keyword": "搜索关键词",
  "page": 1,
  "page_size": 20
}
```

- **响应示例**：与文件和文件夹列表相同

#### 用户文件夹创建

- **接口URL**：`/user/user/folder/create`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "name": "文件夹名称",
  "parent_id": "父文件夹ID，根目录为0"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件夹创建成功",
  "data": {
    "id": "新创建的文件夹ID"
  }
}
```

#### 用户文件夹列表

- **接口URL**：`/user/user/folder/list`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "parent_id": "父文件夹ID，根目录为0"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": "文件夹ID",
        "name": "文件夹名称",
        "parent_id": "父文件夹ID",
        "created_at": "创建时间",
        "updated_at": "更新时间"
      }
    ]
  }
}
```

#### 获取用户文件夹总大小

- **接口URL**：`/user/user/folder/size`
- **请求方式**：GET
- **认证要求**：需要认证
- **请求参数**：无
- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total_size": "总大小(字节)",
    "used_size": "已使用(字节)",
    "available_size": "可用(字节)"
  }
}
```

#### 用户文件的关联存储(文件与文件夹)

- **接口URL**：`/user/user/repository/save`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "name": "文件名",
  "ext": "文件扩展名",
  "size": "文件大小(字节)",
  "path": "文件路径",
  "parent_id": "父文件夹ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件关联成功",
  "data": {
    "id": "关联ID"
  }
}
```

### 回收站管理接口

#### 用户回收站文件删除

- **接口URL**：`/user/user/recycle/delete`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "recycle_id": "回收站项目ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "永久删除成功"
}
```

#### 用户回收站列表

- **接口URL**：`/user/user/recycle/list`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "page": 1,
  "page_size": 20
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total": 5,
    "list": [
      {
        "id": "回收项ID",
        "name": "文件名/文件夹名",
        "ext": "文件扩展名",
        "size": "文件大小(字节)",
        "type": "类型(file/folder)",
        "deleted_at": "删除时间",
        "expires_at": "过期时间(自动永久删除)"
      }
    ]
  }
}
```

#### 用户回收站文件恢复

- **接口URL**：`/user/user/recycle/restore`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "recycle_id": "回收站项目ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "文件恢复成功"
}
```

## 上传服务

上传服务负责处理文件上传、下载和历史记录管理。

### 文件上传下载接口

#### 获取文件下载链接

- **接口URL**：`/upload/file/download/url`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "file_id": "文件ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "download_url": "文件下载URL",
    "expires_in": "URL过期时间(秒)"
  }
}
```

#### 普通文件上传

- **接口URL**：`/upload/file/upload`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求格式**：multipart/form-data
- **请求参数**：
  - **file**: 文件
  - **path**: 保存路径
  - **parent_id**: 父文件夹ID

- **响应示例**：

```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "file_id": "文件ID",
    "name": "文件名",
    "size": "文件大小",
    "path": "文件路径",
    "ext": "文件扩展名"
  }
}
```

### 分片上传接口

#### 初始化分片上传

- **接口URL**：`/upload/file/multipart/init`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "filename": "文件名",
  "size": "文件总大小(字节)",
  "parent_id": "父文件夹ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "初始化成功",
  "data": {
    "upload_id": "上传ID",
    "key": "文件标识"
  }
}
```

#### 上传分片

- **接口URL**：`/upload/file/multipart/upload`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求格式**：multipart/form-data
- **请求参数**：
  - **chunk**: 分片数据
  - **upload_id**: 上传ID
  - **part_number**: 分片序号
  - **key**: 文件标识

- **响应示例**：

```json
{
  "code": 200,
  "message": "分片上传成功",
  "data": {
    "etag": "分片标识",
    "part_number": "分片序号"
  }
}
```

#### 查询分片上传状态

- **接口URL**：`/upload/file/multipart/status`
- **请求方式**：GET
- **认证要求**：需要认证
- **请求参数**：
  - **upload_id**: 上传ID
  - **key**: 文件标识

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "parts": [
      {
        "part_number": "分片序号",
        "etag": "分片标识",
        "size": "分片大小"
      }
    ]
  }
}
```

#### 完成分片上传

- **接口URL**：`/upload/file/multipart/complete`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "upload_id": "上传ID",
  "key": "文件标识",
  "parts": [
    {
      "part_number": "分片序号",
      "etag": "分片标识"
    }
  ]
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "上传完成",
  "data": {
    "file_id": "文件ID",
    "path": "文件路径"
  }
}
```

### 历史记录管理接口

#### 分页查询历史记录列表

- **接口URL**：`/upload/file/history/list`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "page": 1,
  "page_size": 20
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total": 10,
    "list": [
      {
        "id": "历史记录ID",
        "name": "文件名",
        "ext": "文件扩展名",
        "size": "文件大小",
        "path": "文件路径",
        "status": "状态(0:上传中,1:已完成,2:失败)",
        "created_at": "创建时间",
        "updated_at": "更新时间"
      }
    ]
  }
}
```

#### 更新历史上传记录

- **接口URL**：`/upload/file/history/update`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "history_id": "历史记录ID",
  "status": "状态(0:上传中,1:已完成,2:失败)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### 删除所有历史记录

- **接口URL**：`/upload/file/history/delete/all`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "history_ids": ["历史记录ID1", "历史记录ID2"]
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "删除成功"
}
```

## 分享服务

分享服务提供文件分享和协作功能。

### 分享管理接口

#### 创建分享记录

- **接口URL**：`/share/share/basic/create`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "file_id": "文件ID",
  "expire_time": "过期时间(秒)，0表示永不过期",
  "password": "访问密码(可选)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "分享创建成功",
  "data": {
    "share_id": "分享ID",
    "share_url": "分享链接"
  }
}
```

#### 获取资源详情(用于打开分享链接)

- **接口URL**：`/share/share/basic/detail`
- **请求方式**：GET
- **认证要求**：无需认证
- **请求参数**：
  - **share_id**: 分享ID
  - **password**: 访问密码(如果有设置)

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "file": {
      "id": "文件ID",
      "name": "文件名",
      "ext": "文件扩展名",
      "size": "文件大小",
      "type": "类型(file/folder)"
    },
    "creator": {
      "username": "分享创建者用户名",
      "avatar": "头像URL"
    },
    "created_at": "创建时间",
    "expire_time": "过期时间",
    "downloaded": "下载次数"
  }
}
```

#### 资源删除

- **接口URL**：`/share/share/basic/delete`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "share_id": "分享ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "分享删除成功"
}
```

#### 用户分享列表

- **接口URL**：`/share/share/basic/list`
- **请求方式**：GET
- **认证要求**：需要认证
- **请求参数**：
  - **page**: 页码(可选)
  - **page_size**: 每页数量(可选)

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total": 5,
    "list": [
      {
        "id": "分享ID",
        "file_name": "文件名",
        "share_url": "分享链接",
        "expire_time": "过期时间",
        "created_at": "创建时间",
        "downloaded": "下载次数",
        "has_password": "是否有密码(true/false)"
      }
    ]
  }
}
```

#### 资源保存(保存到自己的云盘)

- **接口URL**：`/share/share/basic/save`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "share_id": "分享ID",
  "parent_id": "目标文件夹ID",
  "password": "访问密码(如果有设置)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "file_id": "新文件ID"
  }
}
```

## 日志服务

日志服务负责记录系统操作日志、监控系统状态和SSH连接管理。

### 操作日志接口

#### 分页获得操作日志

- **接口URL**：`/log/operation`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "page": 1,
  "page_size": 20,
  "start_time": "开始时间(可选)",
  "end_time": "结束时间(可选)",
  "operation_type": "操作类型(可选)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "total": 100,
    "list": [
      {
        "id": "日志ID",
        "user_id": "用户ID",
        "username": "用户名",
        "operation": "操作描述",
        "operation_type": "操作类型",
        "ip": "IP地址",
        "user_agent": "浏览器信息",
        "created_at": "创建时间"
      }
    ]
  }
}
```

### 监控接口

#### 历史分析

- **接口URL**：`/log/monitor/history`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "start_time": "开始时间",
  "end_time": "结束时间",
  "metric_type": "指标类型(cpu/memory/disk/network)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "metrics": [
      {
        "timestamp": "时间戳",
        "value": "指标值"
      }
    ],
    "summary": {
      "avg": "平均值",
      "max": "最大值",
      "min": "最小值"
    }
  }
}
```

#### 实时监控

- **接口URL**：`/log/monitor/realtime`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "metric_types": ["cpu", "memory", "disk", "network"]
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "cpu": {
      "usage": "CPU使用率",
      "cores": "CPU核心数",
      "load": "负载"
    },
    "memory": {
      "total": "总内存",
      "used": "已使用",
      "free": "可用",
      "usage_percent": "使用百分比"
    },
    "disk": {
      "total": "总空间",
      "used": "已使用",
      "free": "可用",
      "usage_percent": "使用百分比"
    },
    "network": {
      "rx_bytes": "接收字节",
      "tx_bytes": "发送字节",
      "rx_packets": "接收包数",
      "tx_packets": "发送包数"
    }
  }
}
```

### SSH连接管理接口

#### SSH连接

- **接口URL**：`/log/ssh/connect`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "host": "主机地址",
  "port": "端口号",
  "username": "用户名",
  "password": "密码",
  "description": "连接描述"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "连接成功",
  "data": {
    "connection_id": "连接ID"
  }
}
```

#### 删除SSH连接信息

- **接口URL**：`/log/ssh/delete`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "connection_id": "连接ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "删除成功"
}
```

#### 获取SSH连接信息

- **接口URL**：`/log/ssh/get`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "connection_id": "连接ID"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": "连接ID",
    "host": "主机地址",
    "port": "端口号",
    "username": "用户名",
    "description": "连接描述",
    "created_at": "创建时间"
  }
}
```

#### 获取日志文件列表

- **接口URL**：`/log/ssh/readlog`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "connection_id": "连接ID",
  "path": "目录路径(可选)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "files": [
      {
        "name": "文件名",
        "size": "文件大小",
        "is_dir": "是否为目录",
        "modified_time": "修改时间"
      }
    ]
  }
}
```

#### 读取日志文件

- **接口URL**：`/log/ssh/logfiles`
- **请求方式**：POST
- **认证要求**：需要认证
- **请求参数**：

```json
{
  "connection_id": "连接ID",
  "file_path": "文件路径",
  "lines": "读取行数(可选，默认100)"
}
```

- **响应示例**：

```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "content": "文件内容",
    "total_lines": "总行数"
  }
}
```

## 错误码说明

| 错误码 | 说明             |
| ------ | ---------------- |
| 200    | 成功             |
| 400    | 请求参数错误     |
| 401    | 未授权或授权过期 |
| 403    | 权限不足         |
| 404    | 资源不存在       |
| 409    | 资源冲突         |
| 413    | 文件大小超过限制 |
| 429    | 请求过于频繁     |
| 500    | 服务器内部错误   |
| 503    | 服务不可用       |