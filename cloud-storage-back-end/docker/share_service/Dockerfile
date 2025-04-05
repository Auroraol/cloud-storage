FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译分享服务
RUN go build -o share_service_api ./share_service/api/shareservice.go

FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 复制配置文件和编译好的二进制文件
COPY --from=builder /app/share_service/api/etc /app/etc
COPY --from=builder /app/share_service_api /app/
COPY --from=builder /app/.env /app/

# 创建日志目录
RUN mkdir -p /app/logs

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露API端口
EXPOSE 1006

# 运行服务
CMD ["./share_service_api", "-f", "./etc/shareservice.yaml"] 