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

## 服务访问

部署完成后，可通过以下地址访问各服务：

### 原有项目
- 网站首页：http://101.37.165.220/
- 原有API：http://101.37.165.220/prod-api/

### 云存储系统API
- 用户服务：http://101.37.165.220/api/user/
- 上传服务：http://101.37.165.220/api/upload/
- 分享服务：http://101.37.165.220/api/share/
- 日志服务：http://101.37.165.220/api/log/

## 服务启动顺序

### 第一阶段：RPC 服务启动
1. **上传服务 RPC**
2. **日志服务 RPC**
3. **用户中心 RPC**

### 第二阶段：API 服务启动
1. **用户中心 API**
2. **上传服务 API** 
3. **日志服务 API**
4. **分享服务 API**