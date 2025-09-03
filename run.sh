#!/bin/bash

# 股票指数预测平台启动脚本

echo "🚀 启动股票指数预测平台..."

# 检查Python是否安装
if ! command -v python3 &> /dev/null; then
    echo "❌ 错误: Python3 未安装"
    echo "请先安装Python3: https://www.python.org/downloads/"
    exit 1
fi

# 检查是否在虚拟环境中
if [[ "$VIRTUAL_ENV" == "" ]]; then
    echo "⚠️  警告: 建议在虚拟环境中运行"
    echo "创建虚拟环境: python3 -m venv venv"
    echo "激活虚拟环境: source venv/bin/activate"
    echo ""
fi

# 检查依赖是否安装
echo "📦 检查依赖..."
if ! python3 -c "import flask" 2>/dev/null; then
    echo "📥 安装依赖中..."
    pip3 install -r requirements.txt
    if [ $? -ne 0 ]; then
        echo "❌ 依赖安装失败"
        exit 1
    fi
fi

# 启动服务器
echo "🌐 启动服务器..."
echo "访问地址: http://localhost:9000"
echo "按 Ctrl+C 停止服务器"
echo ""

python3 app.py
