#!/bin/bash

# gogotou - Docker开发环境脚本

echo "🐳 gogotou - Docker开发环境"
echo "=========================================="

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "❌ 错误: Docker daemon未运行"
    echo "请启动Docker Desktop"
    exit 1
fi

# 停止现有开发容器
echo "🛑 停止现有开发容器..."
docker-compose -f docker-compose.dev.yml down

# 构建开发镜像
echo "🔨 构建开发镜像..."
docker-compose -f docker-compose.dev.yml build

# 启动开发服务
echo "🚀 启动开发服务（支持代码热重载）..."
docker-compose -f docker-compose.dev.yml up -d

# 检查服务状态
echo "⏳ 等待服务启动..."
sleep 5

if curl -s http://localhost:9000/api/status > /dev/null; then
    echo "✅ 开发环境启动成功!"
    echo ""
    echo "🌐 访问地址: http://localhost:9000"
    echo "📝 代码修改会自动重载"
    echo ""
    echo "📋 开发命令:"
    echo "  查看日志: docker-compose -f docker-compose.dev.yml logs -f"
    echo "  停止服务: docker-compose -f docker-compose.dev.yml down"
    echo "  进入容器: docker exec -it zhitou-prediction-dev bash"
else
    echo "❌ 开发环境启动失败"
    docker-compose -f docker-compose.dev.yml logs
fi
