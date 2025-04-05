#!/bin/bash

# 颜色定义
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

# 检查项目目录
PROJECT_DIR="/root/cloud-storage"
if [ ! -d "$PROJECT_DIR" ]; then
    print_error "项目目录 $PROJECT_DIR 不存在，请先执行部署脚本"
    exit 1
fi

# 进入项目目录
cd $PROJECT_DIR

# 备份当前配置
backup_configs() {
    print_info "备份当前配置文件..."
    
    BACKUP_DIR="$PROJECT_DIR/backup/$(date +%Y%m%d%H%M%S)"
    mkdir -p $BACKUP_DIR
    
    cp docker-compose.yml $BACKUP_DIR/ 2>/dev/null || true
    cp integrated_nginx.conf $BACKUP_DIR/ 2>/dev/null || true
    cp prometheus.yml $BACKUP_DIR/ 2>/dev/null || true
    cp .env $BACKUP_DIR/ 2>/dev/null || true
    
    print_info "配置已备份到 $BACKUP_DIR"
}

# 更新配置文件
update_configs() {
    print_info "更新配置文件..."
    
    # 如果新文件存在，则替换旧文件
    if [ -f "integrated_docker-compose.yml" ]; then
        cp integrated_docker-compose.yml docker-compose.yml
        print_info "已更新 docker-compose.yml 文件"
    fi
    
    if [ -f "integrated_nginx.conf" ]; then
        cp integrated_nginx.conf /opt/project/nginx/nginx.conf
        print_info "已更新 Nginx 配置文件"
    fi
    
    if [ -f "prometheus.yml" ]; then
        cp prometheus.yml /opt/project/prometheus/prometheus.yml
        print_info "已更新 Prometheus 配置文件"
    fi
}

# 重启服务
restart_services() {
    print_info "重启服务..."
    
    docker-compose down
    docker-compose up -d
    
    if [ $? -ne 0 ]; then
        print_error "服务启动失败，正在回滚..."
        # 如果有备份，尝试恢复
        if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR" ]; then
            cp $BACKUP_DIR/docker-compose.yml . 2>/dev/null || true
            cp $BACKUP_DIR/integrated_nginx.conf /opt/project/nginx/nginx.conf 2>/dev/null || true
            cp $BACKUP_DIR/prometheus.yml /opt/project/prometheus/prometheus.yml 2>/dev/null || true
            cp $BACKUP_DIR/.env . 2>/dev/null || true
            
            docker-compose up -d
            print_warning "已回滚到之前的配置"
        else
            print_error "无法回滚，备份不存在"
        fi
        exit 1
    fi
    
    print_info "服务已成功重启"
}

# 显示服务状态
show_status() {
    print_info "当前服务状态:"
    docker-compose ps
}

# 主函数
main() {
    print_info "开始更新云存储系统..."
    
    backup_configs
    update_configs
    restart_services
    show_status
    
    print_info "服务更新完成，您可以通过以下地址访问服务:"
    echo -e "${GREEN}原有项目:${NC} http://101.37.165.220/"
    echo -e "${GREEN}原有API:${NC} http://101.37.165.220/prod-api/"
    echo -e "${GREEN}------------------------${NC}"
    echo -e "${GREEN}新增API端点:${NC}"
    echo -e "${GREEN}· 用户服务:${NC} http://101.37.165.220/api/user/"
    echo -e "${GREEN}· 上传服务:${NC} http://101.37.165.220/api/upload/"
    echo -e "${GREEN}· 分享服务:${NC} http://101.37.165.220/api/share/"
    echo -e "${GREEN}· 日志服务:${NC} http://101.37.165.220/api/log/"
    echo -e "${GREEN}------------------------${NC}"
    echo -e "${GREEN}监控服务:${NC}"
    echo -e "${GREEN}· Prometheus:${NC} http://101.37.165.220/monitoring/prometheus/"
    echo -e "${GREEN}· Grafana:${NC} http://101.37.165.220/monitoring/grafana/ (用户名: admin, 密码: admin)"
}

# 执行主函数
main 