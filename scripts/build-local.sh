#!/bin/bash

# 智投预测 - 本地构建脚本
# 用于在Docker网络有问题时使用本地Go环境构建

set -e

echo "🔨 使用本地Go环境构建后端..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go 1.21+"
    exit 1
fi

# 进入后端目录
cd backend-go

# 设置Go代理
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 下载依赖
echo "📦 下载Go依赖..."
go mod download

# 构建应用
echo "🔨 构建Go应用..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o ../main ./cmd/main.go

# 返回项目根目录
cd ..

# 创建简化的Dockerfile
cat > Dockerfile.local << 'EOF'
FROM alpine:latest

# 安装 ca-certificates 和 curl
RUN apk --no-cache add ca-certificates curl tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /root/

# 复制构建好的二进制文件
COPY main .

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
EOF

# 构建Docker镜像
echo "🐳 构建Docker镜像..."
docker build -f Dockerfile.local -t alanwzliang/zhitou-prediction-backend:latest .

# 清理临时文件
rm -f main Dockerfile.local

echo "✅ 构建完成！"
echo "🚀 可以使用以下命令启动服务："
echo "   docker run -d -p 8000:8000 --name zhitou-backend alanwzliang/zhitou-prediction-backend:latest"
