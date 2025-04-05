# 云存储系统集成部署指南

## 项目概述

本项目是一个云存储系统与现有项目的集成部署方案，旨在实现以下目标：

1. 在不影响现有项目功能的前提下，部署新的云存储系统
2. 提供统一的访问入口，使两个项目能够共存
3. 建立完整的监控体系，确保系统稳定运行

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
- **监控系统**：
  - Prometheus：指标收集与存储
  - Grafana：数据可视化面板

## 集成方案

为了确保现有项目和新系统的兼容性，我们采用以下策略：

1. **统一的Nginx代理**：使用集成的Nginx配置，通过不同路径区分两个项目的请求
2. **隔离的网络环境**：所有服务通过Docker网络互相通信，确保安全性
3. **监控系统整合**：为所有服务提供监控能力，便于问题排查

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

### 监控服务
- Prometheus：http://101.37.165.220/monitoring/prometheus/
- Grafana：http://101.37.165.220/monitoring/grafana/（默认用户名/密码：admin/admin）

## 常见问题

### 如何检查服务状态？

```bash
cd /root/cloud-storage
docker-compose ps
```

### 如何查看服务日志？

```bash
# 查看特定服务的日志
docker-compose logs [服务名]

# 例如查看用户中心服务的日志
docker-compose logs user-center
```

### 如何重启特定服务？

```bash
docker-compose restart [服务名]
```

### 如何解决端口冲突？

如果遇到端口冲突，可以修改`integrated_docker-compose.yml`文件中的端口映射，然后执行更新脚本。

### 如何备份数据？

```bash
# 备份MySQL数据
docker exec mysql mysqldump -u root -p123456 cloud_storage > backup_$(date +%Y%m%d).sql
```

## 安全建议

1. 修改默认密码：部署后请立即修改MySQL、Redis和Grafana的默认密码
2. 启用HTTPS：建议配置SSL证书，启用HTTPS访问
3. 限制端口访问：只对外开放必要的端口，如80和443
4. 定期备份：建立定期备份机制，确保数据安全

## 维护指南

### 日常维护

- 定期检查磁盘空间使用情况
- 监控服务资源使用情况
- 查看Grafana监控面板，关注系统健康状态

### 故障排除

1. 服务无法访问：
   - 检查Docker容器是否运行：`docker-compose ps`
   - 检查Nginx配置是否正确：`docker exec nginx nginx -t`
   - 检查防火墙设置：`firewall-cmd --list-all`

2. 数据库连接问题：
   - 检查MySQL容器状态：`docker-compose logs mysql`
   - 验证数据库连接：`docker exec -it mysql mysql -u root -p123456`

3. 监控系统问题：
   - 检查Prometheus配置：`cat /opt/project/prometheus/prometheus.yml`
   - 确认指标采集目标可达：访问Prometheus目标页面查看状态

## 联系与支持

如遇到部署或使用问题，请联系技术支持团队：

- 邮箱：support@example.com
- 电话：123-456-7890

---

© 2023 云存储项目团队 