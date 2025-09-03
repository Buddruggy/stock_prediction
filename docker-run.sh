#!/bin/bash

# 智投预测 - Docker部署脚本

echo "🐳 智投预测 - AI股市指数预测平台"
echo "=========================================="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker未安装"
    echo "请先安装Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查Docker daemon是否运行
if ! docker info &> /dev/null; then
    echo "❌ 错误: Docker daemon未运行"
    echo "请启动Docker Desktop或运行以下命令:"
    echo "  macOS: open -a Docker"
    echo "  Linux: sudo systemctl start docker"
    echo "  或使用Colima: colima start"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: Docker Compose未安装"
    echo "请先安装Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# 停止并移除现有容器
echo "🛑 停止现有容器..."
docker-compose down --remove-orphans

# 构建镜像
echo "🔨 构建Docker镜像..."
docker-compose build --no-cache

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

# 检查服务状态
echo "⏳ 等待服务启动..."
sleep 10

# 健康检查
echo "🔍 检查服务健康状态..."
if curl -s http://localhost:9000/api/status > /dev/null; then
    echo "✅ 服务启动成功!"
    echo ""
    echo "🌐 访问地址: http://localhost:9000"
    echo "📊 API状态: http://localhost:9000/api/status"
    echo ""
    echo "📋 Docker命令:"
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
    echo "  重启服务: docker-compose restart"
    echo ""
else
    echo "❌ 服务启动失败，请检查日志:"
    echo "docker-compose logs"
fi
