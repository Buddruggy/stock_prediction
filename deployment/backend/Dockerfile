# 智投预测 - AI股市指数预测平台 Docker镜像 (Go后端)
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 设置Go代理和模块配置（使用国内镜像源）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 复制 go mod 文件
COPY backend-go/go.mod backend-go/go.sum ./

# 下载依赖（使用详细输出和重试机制）
RUN go mod download -x || (sleep 5 && go mod download -x) || (sleep 10 && go mod download -x)

# 复制源代码
COPY backend-go/ .

# 构建应用
RUN go build -a -installsuffix cgo -ldflags '-w -s' -o main ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 和 curl
RUN apk --no-cache add ca-certificates curl tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建非root用户
RUN adduser -D -s /bin/sh app && chown -R app:app /root
USER app

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8000/health || exit 1

# 暴露端口
EXPOSE 8000

# 运行应用
CMD ["./main"]