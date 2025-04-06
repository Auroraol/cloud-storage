# frp内网穿透

## 1. 什么是内网穿透

内网穿透（Port Forwarding）是一种网络技术，允许你将外部网络中的请求转发到内部网络中的特定计算机或设备。

内网穿透通常需要一种中间代理或服务器，它位于内部网络和外部网络之间，将外部请求路由到内部目标。这可以通过各种工具和服务来实现，包括专用的内网穿透工具、虚拟专用网络（VPN）和云服务。

其中，frp（Fast Remote Port Forwarding）是一款流行的内网穿透工具。

<img src="frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250407004429271.png" alt="image-20250407004429271" style="zoom:67%;" />

## 2. 实现背景

go微服务后端部署难度过大.  所以最后实现的将是：在任意一台计算上通过访问公网ip就能访问我本地计算机的服务

## 3. 配置文档



## 4. 实现

### 服务器端（云服务器）

服务端选 Linux 安装包：frp_0.61.1_linux_amd64.tar.gz

![image-20250407003512939](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250407003512939.png)

#### 云服务器配置

![image-20250406222852286](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250406222852286.png)

#### frps.toml

```shell
cd /opt/frps
vim frps.toml
```

frps.tonl

```
bindPort = 7000

# 客户端连接服务端的认证
auth.token = "lfj@1665834"

# web界面配置
webServer.addr = "0.0.0.0"
webServer.port = 7500
webServer.user = "admin"
webServer.password = "admin"
```

#### docker-compose.yml

```shell
cd /opt/frps
vim docker-compose.yml
```

docker-compose.yml

```dockerfile
version: "3.8"
services:
  frps:
    image: snowdreamtech/frps:0.61
    container_name: frps
    ports:
      - 7000:7000
      - 7500:7500
    volumes:
      - /opt/frps/frps.toml:/etc/frp/frps.toml
      - /opt/frps/logs.log:/etc/frp/frps.log
```

启动

```shell
docker-compose up --build -d
docker exec -it frps sh
```

#### 效果

docker:

![image-20250407003039233](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250407003039233.png)

frps dashboard:  http://101.37.165.220:7500

![image-20250406222211586](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250406222211586.png)

### 客户端（本地电脑）

客户端选 Windows 安装包：frp_0.61.1_windows_amd64

![image-20250407003734032](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250407003734032.png)

#### frpc.toml

```
serverAddr = "101.37.165.220" # 阿里云IP
serverPort = 7000
auth.token = "lfj@1665834"
[[proxies]]
name = "http_user_center"
type = "tcp"
localIP = "127.0.0.1"
localPort = 1004
remotePort = 1004

[[proxies]]
name = "http_file_upload"
type = "tcp"
localIP = "127.0.0.1"
localPort = 1005
remotePort  = 1005

[[proxies]]
name = "http_share"
type = "tcp"
localIP = "127.0.0.1"
localPort = 1006
remotePort  = 1006

[[proxies]]
name = "http_log_service"
type = "tcp"
localIP = "127.0.0.1"
localPort = 1007
remotePort  = 1007
```

#### 运行

```shell
frpc -c ./frpc.toml
```

#### 效果

![image-20250407004038297](frp%E5%86%85%E7%BD%91%E7%A9%BF%E9%80%8F.assets/image-20250407004038297.png)
