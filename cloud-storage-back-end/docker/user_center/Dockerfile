FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译用户中心服务
RUN go build -o user_center_api ./user_center/api/usercenter.go
RUN go build -o user_center_rpc ./user_center/rpc/usercenter.go

FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 复制配置文件和编译好的二进制文件
COPY --from=builder /app/user_center/api/etc /app/etc
COPY --from=builder /app/user_center_api /app/
COPY --from=builder /app/user_center_rpc /app/
COPY --from=builder /app/.env /app/

# 创建日志目录
RUN mkdir -p /app/logs

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露API端口
EXPOSE 1004

# 运行服务
CMD ["sh", "-c", "./user_center_api -f ./etc/usercenter.yaml & ./user_center_rpc -f ./etc/usercenter-rpc.yaml"] 