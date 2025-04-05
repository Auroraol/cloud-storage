FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译文件处理服务
RUN go build -o file_processor ./file_processor/main.go

FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 复制配置文件和编译好的二进制文件
COPY --from=builder /app/file_processor/etc /app/etc
COPY --from=builder /app/file_processor /app/
COPY --from=builder /app/.env /app/

# 创建日志目录
RUN mkdir -p /app/logs

# 设置时区
ENV TZ=Asia/Shanghai

# 运行服务
CMD ["./file_processor", "-f", "./etc/config.yaml"] 