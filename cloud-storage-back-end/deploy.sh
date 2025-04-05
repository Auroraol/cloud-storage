#!/bin/bash

# 显示颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 输出带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker未安装，请先安装Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
}

# 检查运行环境
check_environment() {
    # 检查现有后端服务是否运行
    if ! curl -s http://localhost:9090 &>/dev/null; then
        print_warning "端口9090上没有检测到服务，请确保您的后端服务正在运行"
        read -p "是否继续部署? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        print_info "已检测到9090端口后端服务正在运行"
    fi

    # 检查操作系统是否为Linux，如果是需要设置host.docker.internal
    if [[ "$(uname)" == "Linux" ]]; then
        print_info "Linux系统检测到，设置host.docker.internal域名"
        
        # 获取Docker网桥网关IP
        DOCKER_GATEWAY=$(docker network inspect bridge | grep Gateway | awk '{print $2}' | tr -d '"')
        
        if [ -z "$DOCKER_GATEWAY" ]; then
            print_warning "无法获取Docker网关IP，使用默认值172.17.0.1"
            DOCKER_GATEWAY="172.17.0.1"
        fi
        
        # 检查/etc/hosts中是否已存在host.docker.internal
        if grep -q "host.docker.internal" /etc/hosts; then
            print_info "host.docker.internal已存在于/etc/hosts文件中"
        else
            print_info "将host.docker.internal添加到/etc/hosts文件中"
            echo "$DOCKER_GATEWAY host.docker.internal" | sudo tee -a /etc/hosts
            if [ $? -ne 0 ]; then
                print_error "无法添加host.docker.internal到/etc/hosts，请手动添加或以管理员身份运行"
                print_info "请手动执行: echo \"$DOCKER_GATEWAY host.docker.internal\" | sudo tee -a /etc/hosts"
            fi
        fi
    fi
}

# 创建必要的目录
create_directories() {
    print_info "创建必要的目录..."
    
    mkdir -p /opt/project/mysql/data
    mkdir -p /opt/project/mysql/conf
    mkdir -p /opt/project/redis/data
    mkdir -p /opt/project/redis/conf
    mkdir -p /opt/project/pulsar_data
    mkdir -p /opt/project/etcd_data
    mkdir -p /opt/project/nginx/html
    mkdir -p /opt/project/nginx/log
    mkdir -p /opt/project/ssl
    
    # 添加Prometheus和Grafana目录
    mkdir -p /opt/project/prometheus/data
    mkdir -p /opt/project/grafana/data
    
    # 添加应用服务日志目录
    mkdir -p /opt/project/logs
    
    # 创建一个临时的Redis配置文件，如果没有的话
    if [ ! -f "/opt/project/redis/conf/redis.conf" ]; then
        print_warning "未找到Redis配置文件，创建默认配置..."
        echo "# Redis配置文件
port 6379
protected-mode yes
dir /data
appendonly yes" > /opt/project/redis/conf/redis.conf
    fi
    
    # 复制Prometheus配置文件
    if [ -f "prometheus.yml" ]; then
        cp prometheus.yml /opt/project/prometheus/prometheus.yml
        print_info "已复制Prometheus配置文件"
    else
        print_warning "未找到prometheus.yml文件，请确保该文件存在"
    fi
}

# 复制Nginx配置文件
setup_nginx() {
    print_info "设置Nginx配置..."
    
    cp nginx.conf /opt/project/nginx/nginx.conf
    if [ $? -ne 0 ]; then
        print_error "复制Nginx配置文件失败"
        exit 1
    fi
}

# 启动服务
start_services() {
    print_info "开始部署服务..."
    
    # 检查.env文件是否存在
    if [ ! -f ".env" ]; then
        print_warning "未找到.env文件，从.env.example创建..."
        cp .env.example .env
        print_info "请记得编辑.env文件，更新您的阿里云OSS配置"
    fi
    
    docker-compose up -d
    if [ $? -ne 0 ]; then
        print_error "启动服务失败，请检查日志"
        exit 1
    fi
    
    print_info "所有服务已成功启动"
}

# 显示服务状态
show_status() {
    print_info "当前服务状态:"
    docker-compose ps
}

# 主函数
main() {
    print_info "开始部署云存储系统..."
    
    check_docker
    check_environment
    create_directories
    setup_nginx
    start_services
    show_status
    
    print_info "部署完成，您可以通过以下地址访问系统:"
    echo -e "${GREEN}API网关:${NC} http://101.37.165.220/"
    echo -e "${GREEN}Pulsar管理界面:${NC} http://101.37.165.220:8080/"
    echo -e "${GREEN}Prometheus监控:${NC} http://101.37.165.220:9090/"
    echo -e "${GREEN}Grafana仪表盘:${NC} http://101.37.165.220:3000/ (用户名: admin, 密码: admin)"
    
    print_info "以下API端点已配置:"
    echo -e "${GREEN}用户服务:${NC} http://101.37.165.220/api/user/"
    echo -e "${GREEN}上传服务:${NC} http://101.37.165.220/api/upload/"
    echo -e "${GREEN}分享服务:${NC} http://101.37.165.220/api/share/"
    echo -e "${GREEN}日志服务:${NC} http://101.37.165.220/api/log/"
    echo -e "${GREEN}原有后端服务:${NC} http://101.37.165.220/prod-api/"
}

# 执行主函数
main 