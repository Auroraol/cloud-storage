# 云存储系统部署指南

## 服务器准备

1. 登录到您的云服务器:

```bash
ssh root@101.37.165.220
# 输入您的密码
```

2. 安装Docker和Docker Compose:

```bash
# 安装Docker
curl -fsSL https://get.docker.com | bash -s docker

# 安装Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.20.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

## 部署步骤

1. 创建项目目录并进入:

```bash
mkdir -p /root/cloud-storage
cd /root/cloud-storage
```

2. 下载部署文件:

```bash
# 创建必要文件
touch docker-compose.yml nginx.conf prometheus.yml deploy.sh
chmod +x deploy.sh
```

3. 编辑Docker Compose文件:

```bash
vi docker-compose.yml
# 粘贴我们之前准备好的docker-compose.yml内容
```

4. 编辑Nginx配置文件:

```bash
vi nginx.conf
# 粘贴我们之前准备好的nginx.conf内容
```

5. 编辑Prometheus配置文件:

```bash
vi prometheus.yml
# 粘贴我们之前准备好的prometheus.yml内容
```

6. 编辑部署脚本:

```bash
vi deploy.sh
# 粘贴我们之前准备好的deploy.sh内容
```

7. 创建.env文件:

```bash
vi .env
# 添加以下内容并替换为您的实际值
OSS_ACCESS_KEY_ID="您的ACCESS_KEY_ID"
OSS_ACCESS_KEY_SECRET="您的ACCESS_KEY_SECRET"
OSS_BUCKET_NAME="您的BUCKET_NAME"
OSS_ENDPOINT="oss-cn-beijing.aliyuncs.com"
```

8. 执行部署脚本:

```bash
./deploy.sh
```

## 访问服务

部署完成后，您可以通过以下地址访问各服务:

- API网关: http://101.37.165.220/
- Pulsar管理界面: http://101.37.165.220:8080/
- Prometheus监控: http://101.37.165.220:9090/
- Grafana仪表盘: http://101.37.165.220:3000/ (用户名: admin, 密码: admin)

## 安全建议

1. 配置防火墙, 仅开放必要端口:

```bash
# 安装防火墙(如果未安装)
apt-get update && apt-get install -y ufw

# 配置防火墙规则
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 8080/tcp
ufw allow 9090/tcp
ufw allow 3000/tcp
ufw enable
```

2. 设置强密码和SSH密钥登录

```bash
# 生成新的SSH密钥对(在本地计算机上执行)
ssh-keygen -t ed25519 -C "your_email@example.com"

# 将公钥复制到服务器(在本地计算机上执行)
ssh-copy-id root@101.37.165.220

# 在服务器上禁用密码登录
vi /etc/ssh/sshd_config
# 修改以下设置
# PasswordAuthentication no
# 保存退出

# 重启SSH服务
systemctl restart sshd
```

## 维护操作

1. 查看容器运行状态:

```bash
cd /root/cloud-storage
docker-compose ps
```

2. 查看服务日志:

```bash
docker-compose logs -f [服务名]
```

3. 重启特定服务:

```bash
docker-compose restart [服务名]
```

4. 更新整个系统:

```bash
cd /root/cloud-storage
git pull # 如果是从git仓库克隆的
docker-compose down
docker-compose up -d
```