# Go国内代理配置指南

## 🚀 已完成的配置

### 1. Go环境代理设置
```bash
# 设置Go代理（已执行）
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 验证配置
go env GOPROXY
go env GOSUMDB
```

### 2. Dockerfile代理配置
已在Dockerfile中添加以下配置：
```dockerfile
# 设置Go代理和模块配置
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on
```

## 🔧 可用的构建命令

### 方式一：Docker构建（推荐）
```bash
# 等待Docker启动后执行
make build-backend
```

### 方式二：本地构建
```bash
# 使用本地Go环境构建
make build-backend-local
```

### 方式三：手动构建
```bash
# 1. 进入后端目录
cd backend-go

# 2. 下载依赖（使用国内代理）
go mod download

# 3. 构建应用
go build -o main ./cmd/main.go

# 4. 返回根目录
cd ..

# 5. 构建Docker镜像
docker build -f deployment/backend/Dockerfile -t alanwzliang/zhitou-prediction-backend:latest .
```

## 🌐 国内代理源说明

### Go代理源
- **goproxy.cn**: 七牛云提供的Go模块代理
- **direct**: 当代理不可用时直接访问源
- **sum.golang.google.cn**: Google提供的校验和数据库镜像

### 其他可选代理源
```bash
# 阿里云代理
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 腾讯云代理
export GOPROXY=https://mirrors.cloud.tencent.com/goproxy/,direct

# 中科大代理
export GOPROXY=https://goproxy.io,direct
```

## 🚀 部署命令

### 开发环境
```bash
# 直接运行Go服务
make backend
```

### 生产环境
```bash
# 构建并部署
make build-backend
make up
```

### 一键部署
```bash
# 构建并启动所有服务
make up
```

## 📊 服务访问地址

- **后端API**: http://localhost:8000
- **健康检查**: http://localhost:8000/health
- **所有指数**: http://localhost:8000/api/v1/indices/all
- **所有预测**: http://localhost:8000/api/v1/predict/all

## 🔍 故障排除

### Docker连接问题
```bash
# 检查Colima状态
colima status

# 重启Colima
colima stop && colima start

# 检查Docker上下文
docker context ls
docker context use colima
```

### Go模块下载问题
```bash
# 清理模块缓存
go clean -modcache

# 重新下载
go mod download

# 验证代理
go env GOPROXY
```

## ✅ 当前状态

- ✅ Go代理已配置
- ✅ Dockerfile已优化
- ✅ 本地Go环境正常
- ⏳ Docker daemon启动中
- ⏳ 等待Docker构建测试
