# 智投预测 - AI股市指数预测平台 Makefile
# 用于构建、打包和部署Docker镜像

# 项目配置
PROJECT_NAME := zhitou-prediction
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
REGISTRY := docker.io
USERNAME := alanwzliang

# 前后端分离配置
BACKEND_IMAGE := $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend
FRONTEND_IMAGE := $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend
BACKEND_CONTAINER := $(PROJECT_NAME)-backend
FRONTEND_CONTAINER := $(PROJECT_NAME)-frontend

# 端口配置
BACKEND_PORT := 8000
FRONTEND_PORT := 9000
LEGACY_PORT := 9001

# Docker配置
BACKEND_DOCKERFILE := deployment/backend/Dockerfile
FRONTEND_DOCKERFILE := deployment/frontend/Dockerfile
FRONTEND_HTTPS_DOCKERFILE := deployment/frontend/Dockerfile.https
FRONTEND_PROD_DOCKERFILE := deployment/frontend/Dockerfile.prod
LEGACY_DOCKERFILE := deploy/docker/Dockerfile
DOCKER_CONTEXT := .

# 颜色输出
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help build push pull run stop clean logs shell test dev prod deploy backend frontend legacy install-deps build-frontend-https up-https quick-https ssl-cert setup-letsencrypt build-prod up-prod

# 默认目标
.DEFAULT_GOAL := help

## 显示帮助信息
help:
	@echo "$(BLUE)智投预测 - AI股市指数预测平台$(RESET)"
	@echo "$(BLUE)=========================================$(RESET)"
	@echo ""
	@echo "$(GREEN)可用命令:$(RESET)"
	@echo ""
	@echo "$(YELLOW)快速启动:$(RESET)"
	@echo "  up        - 🚀 一键部署前后端服务 (推荐)"
	@echo "  quick     - 🚀 智能选择启动模式"
	@echo "  legacy    - 🚀 单体应用兼容模式"
	@echo ""
	@echo "$(YELLOW)开发相关:$(RESET)"
	@echo "  install-deps - 安装前后端依赖"
	@echo "  backend   - 启动Go后端开发服务器"
	@echo "  frontend  - 启动前端开发服务器"
	@echo "  dev       - 同时启动前后端开发服务器"
	@echo ""
	@echo "$(YELLOW)构建相关:$(RESET)"
	@echo "  build     - 构建前后端Docker镜像"
	@echo "  build-backend - 构建Go后端镜像"
	@echo "  build-frontend - 构建前端镜像"
	@echo "  legacy    - 构建并运行单体应用(兼容)"
	@echo "  push      - 推送镜像到仓库"
	@echo "  pull      - 从仓库拉取镜像"
	@echo ""
	@echo "$(YELLOW)运行相关:$(RESET)"
	@echo "  run       - 运行前后端容器"
	@echo "  prod      - 运行生产环境"
	@echo "  stop      - 停止所有容器"
	@echo "  restart   - 重启容器"
	@echo ""
	@echo "$(YELLOW)管理相关:$(RESET)"
	@echo "  logs      - 查看容器日志"
	@echo "  shell     - 进入容器shell"
	@echo "  ps        - 查看容器状态"
	@echo "  clean     - 清理未使用的镜像"
	@echo "  logs-db   - 查看MySQL日志"
	@echo "  db-shell  - 进入MySQL容器"
	@echo ""
	@echo "$(YELLOW)测试相关:$(RESET)"
	@echo "  test      - 运行测试"
	@echo "  health    - 健康检查"
	@echo ""
	@echo "$(YELLOW)部署相关:$(RESET)"
	@echo "  deploy    - 部署到生产环境"
	@echo "  k8s       - 部署到Kubernetes"
	@echo ""
	@echo "$(GREEN)当前配置:$(RESET)"
	@echo "  项目名称: $(PROJECT_NAME)"
	@echo "  版本号:   $(VERSION)"
	@echo "  后端镜像: $(BACKEND_IMAGE):$(VERSION)"
	@echo "  前端镜像: $(FRONTEND_IMAGE):$(VERSION)"
	@echo "  后端端口: $(BACKEND_PORT)"
	@echo "  前端端口: $(FRONTEND_PORT)"
	@echo "  兼容端口: $(LEGACY_PORT)"

## 安装前后端依赖
install-deps:
	@echo "$(GREEN)📦 安装前后端依赖...$(RESET)"
	@echo "$(YELLOW)安装Go后端依赖...$(RESET)"
	cd backend-go && go mod tidy
	@echo "$(YELLOW)安装前端依赖...$(RESET)"
	cd frontend && npm install
	@echo "$(GREEN)✅ 依赖安装完成$(RESET)"

## 启动Go后端开发服务器
backend:
	@echo "$(GREEN)🚀 启动Go后端开发服务器...$(RESET)"
	@echo "$(BLUE)🌐 后端地址: http://localhost:$(BACKEND_PORT)$(RESET)"
	@echo "$(BLUE)📚 API文档: http://localhost:$(BACKEND_PORT)/api/v1$(RESET)"
	cd backend-go && go run cmd/main.go

## 启动前端开发服务器
frontend:
	@echo "$(GREEN)🚀 启动前端开发服务器...$(RESET)"
	@echo "$(BLUE)🌐 前端地址: http://localhost:$(FRONTEND_PORT)$(RESET)"
	cd frontend && npm run dev

## 同时启动前后端开发服务器
dev:
	@echo "$(GREEN)🚀 启动前后端开发环境...$(RESET)"
	@echo "$(BLUE)🌐 前端: http://localhost:$(FRONTEND_PORT)$(RESET)"
	@echo "$(BLUE)🌐 后端: http://localhost:$(BACKEND_PORT)$(RESET)"
	@echo "$(BLUE)📚 API文档: http://localhost:$(BACKEND_PORT)/api/v1$(RESET)"
	@echo "$(YELLOW)请在两个终端分别运行:$(RESET)"
	@echo "  终端1: make backend"
	@echo "  终端2: make frontend"

## 一键部署前后端服务
up: stop
	@echo "$(GREEN)🚀 一键部署前后端服务...$(RESET)"
	@echo "$(YELLOW)正在构建前后端镜像...$(RESET)"
	@$(MAKE) build-backend build-frontend
	@echo "$(YELLOW)正在创建网络...$(RESET)"
	-docker network create zhitou-network 2>/dev/null || true
	@echo "$(YELLOW)正在启动Go后端服务...$(RESET)"
	docker run -d \
		--name $(BACKEND_CONTAINER) \
		--network zhitou-network \
		-p $(BACKEND_PORT):$(BACKEND_PORT) \
		-e ENVIRONMENT=production \
		--restart unless-stopped \
		$(BACKEND_IMAGE):latest
	@echo "$(YELLOW)正在启动前端服务...$(RESET)"
	docker run -d \
		--name $(FRONTEND_CONTAINER) \
		--network zhitou-network \
		-p 80:80 \
		-p $(FRONTEND_PORT):$(FRONTEND_PORT) \
		--restart unless-stopped \
		$(FRONTEND_IMAGE):latest
	@echo "$(GREEN)✅ 前后端服务部署完成$(RESET)"
	@echo "$(BLUE)🌐 前端访问: http://localhost:80 或 http://localhost:$(FRONTEND_PORT)$(RESET)"
	@echo "$(BLUE)🌐 后端API: http://localhost:$(BACKEND_PORT)$(RESET)"
	@echo "$(BLUE)📚 API文档: http://localhost:$(BACKEND_PORT)/api/v1$(RESET)"
	@echo "$(YELLOW)💡 使用 'make stop' 停止服务$(RESET)"

## 快速一键启动 (推荐)
quick: 
	@echo "$(GREEN)🚀 快速一键启动...$(RESET)"
	@echo "$(YELLOW)选择启动模式:$(RESET)"
	@echo "  1. 前后端分离模式 (推荐)"
	@echo "  2. 兼容单体模式"
	@read -p "请选择 [1/2]: " mode; \
	if [ "$$mode" = "2" ]; then \
		$(MAKE) legacy; \
	else \
		$(MAKE) up; \
	fi

## 构建Go后端Docker镜像
build-backend:
	@echo "$(GREEN)🔨 构建Go后端Docker镜像...$(RESET)"
	docker build -t $(BACKEND_IMAGE):$(VERSION) -t $(BACKEND_IMAGE):latest -f $(BACKEND_DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ Go后端镜像构建完成: $(BACKEND_IMAGE):$(VERSION)$(RESET)"

## 本地构建Go后端Docker镜像（避免网络问题）
build-backend-local:
	@echo "$(GREEN)🔨 使用本地Go环境构建后端镜像...$(RESET)"
	./scripts/build-local.sh
	@echo "$(GREEN)✅ 本地构建完成$(RESET)"

## 构建前端Docker镜像
build-frontend:
	@echo "$(GREEN)🔨 构建前端Docker镜像...$(RESET)"
	docker build -t $(FRONTEND_IMAGE):$(VERSION) -t $(FRONTEND_IMAGE):latest -f $(FRONTEND_DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 前端镜像构建完成: $(FRONTEND_IMAGE):$(VERSION)$(RESET)"

## 构建前后端Docker镜像
build: build-backend build-frontend
	@echo "$(GREEN)✅ 前后端镜像构建完成$(RESET)"

## 兼容单体应用
legacy:
	@echo "$(YELLOW)🔄 运行单体应用(兼容模式)...$(RESET)"
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME):$(VERSION) -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME):latest -f $(LEGACY_DOCKERFILE) $(DOCKER_CONTEXT)
	docker stop $(PROJECT_NAME) 2>/dev/null || true
	docker rm $(PROJECT_NAME) 2>/dev/null || true
	docker run -d --name $(PROJECT_NAME) -p $(LEGACY_PORT):$(LEGACY_PORT) $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME):latest
	@echo "$(GREEN)✅ 单体应用启动完成$(RESET)"
	@echo "$(BLUE)🌐 访问地址: http://localhost:$(LEGACY_PORT)$(RESET)"

## 推送镜像到仓库
push: build
	@echo "$(GREEN)📤 推送镜像到仓库...$(RESET)"
	docker push $(BACKEND_IMAGE):$(VERSION)
	docker push $(BACKEND_IMAGE):latest
	docker push $(FRONTEND_IMAGE):$(VERSION)
	docker push $(FRONTEND_IMAGE):latest
	@echo "$(GREEN)✅ 镜像推送完成$(RESET)"

## 从仓库拉取镜像
pull:
	@echo "$(GREEN)📥 从仓库拉取镜像...$(RESET)"
	docker pull $(BACKEND_IMAGE):latest
	docker pull $(FRONTEND_IMAGE):latest
	@echo "$(GREEN)✅ 镜像拉取完成$(RESET)"

## 运行生产环境
prod: build up
	@echo "$(GREEN)🏭 生产环境部署完成$(RESET)"

## 停止所有容器
stop:
	@echo "$(YELLOW)🛑 停止所有容器...$(RESET)"
	-docker stop $(BACKEND_CONTAINER) $(FRONTEND_CONTAINER) $(PROJECT_NAME) 2>/dev/null || true
	-docker stop $(BACKEND_CONTAINER)-dev $(FRONTEND_CONTAINER)-dev $(PROJECT_NAME)-dev 2>/dev/null || true
	-docker rm $(BACKEND_CONTAINER) $(FRONTEND_CONTAINER) $(PROJECT_NAME) 2>/dev/null || true
	-docker rm $(BACKEND_CONTAINER)-dev $(FRONTEND_CONTAINER)-dev $(PROJECT_NAME)-dev 2>/dev/null || true
	@echo "$(GREEN)✅ 容器已停止$(RESET)"

## 重启容器
restart: stop up

## 查看容器日志
logs:
	@echo "$(BLUE)📋 查看容器日志...$(RESET)"
	@echo "$(YELLOW)后端日志:$(RESET)"
	docker logs -f $(BACKEND_CONTAINER) 2>/dev/null || echo "$(RED)后端容器未运行$(RESET)"
	@echo "$(YELLOW)前端日志:$(RESET)"
	docker logs -f $(FRONTEND_CONTAINER) 2>/dev/null || echo "$(RED)前端容器未运行$(RESET)"

## 进入容器shell
shell:
	@echo "$(BLUE)🐚 进入容器shell...$(RESET)"
	@echo "$(YELLOW)选择容器:$(RESET)"
	@echo "  1. Go后端容器"
	@echo "  2. 前端容器"
	@read -p "请选择 [1/2]: " choice; \
	if [ "$$choice" = "1" ]; then \
		docker exec -it $(BACKEND_CONTAINER) /bin/sh; \
	elif [ "$$choice" = "2" ]; then \
		docker exec -it $(FRONTEND_CONTAINER) /bin/sh; \
	else \
		echo "$(RED)无效选择$(RESET)"; \
	fi

## 查看容器状态
ps:
	@echo "$(BLUE)📊 容器状态:$(RESET)"
	docker ps -a --filter "name=$(PROJECT_NAME)"

## 清理未使用的镜像
clean:
	@echo "$(YELLOW)🧹 清理未使用的Docker资源...$(RESET)"
	docker system prune -f
	docker image prune -f
	@echo "$(GREEN)✅ 清理完成$(RESET)"

## 深度清理
clean-all: stop
	@echo "$(RED)🗑 深度清理Docker资源...$(RESET)"
	-docker rmi $(BACKEND_IMAGE):$(VERSION) 2>/dev/null || true
	-docker rmi $(BACKEND_IMAGE):latest 2>/dev/null || true
	-docker rmi $(FRONTEND_IMAGE):$(VERSION) 2>/dev/null || true
	-docker rmi $(FRONTEND_IMAGE):latest 2>/dev/null || true
	docker system prune -af
	@echo "$(GREEN)✅ 深度清理完成$(RESET)"

## 生成SSL证书
ssl-cert:
	@echo "$(GREEN)🔐 生成SSL证书...$(RESET)"
	@echo "$(YELLOW)请输入您的域名 (默认: localhost):$(RESET)"
	@read -p "域名: " domain; \
	if [ -z "$$domain" ]; then domain="localhost"; fi; \
	./scripts/generate-ssl-cert.sh "$$domain"
	@echo "$(GREEN)✅ SSL证书生成完成$(RESET)"
	@echo "$(BLUE)📁 证书位置: ssl/cert.pem, ssl/key.pem$(RESET)"

## 构建HTTPS前端镜像
build-frontend-https:
	@echo "$(GREEN)🔨 构建HTTPS前端Docker镜像...$(RESET)"
	docker build -t $(FRONTEND_IMAGE):$(VERSION)-https -t $(FRONTEND_IMAGE):latest-https -f $(FRONTEND_HTTPS_DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ HTTPS前端镜像构建完成: $(FRONTEND_IMAGE):$(VERSION)-https$(RESET)"

## 部署HTTPS服务
up-https: build-backend build-frontend-https
	@echo "$(GREEN)🚀 一键部署HTTPS前后端服务...$(RESET)"
	@echo "$(YELLOW)正在构建前后端镜像...$(RESET)"
	@echo "$(YELLOW)正在创建网络...$(RESET)"
	docker network create zhitou-network 2>/dev/null || true
	@echo "$(YELLOW)正在启动Go后端服务...$(RESET)"
	docker run -d \
		--name $(BACKEND_CONTAINER) \
		--network zhitou-network \
		-p $(BACKEND_PORT):$(BACKEND_PORT) \
		-e ENVIRONMENT=production \
		--restart unless-stopped \
		$(BACKEND_IMAGE):latest
	@echo "$(YELLOW)正在启动HTTPS前端服务...$(RESET)"
	docker run -d \
		--name $(FRONTEND_CONTAINER)-https \
		--network zhitou-network \
		-p 80:80 \
		-p 443:443 \
		-v "$(PWD)/ssl:/etc/nginx/ssl:ro" \
		--restart unless-stopped \
		$(FRONTEND_IMAGE):latest-https
	@echo "$(GREEN)✅ HTTPS前后端服务部署完成$(RESET)"
	@echo "$(BLUE)🌐 HTTPS访问: https://localhost:443$(RESET)"
	@echo "$(BLUE)🌐 HTTP重定向: http://localhost:80 -> https://localhost:443$(RESET)"
	@echo "$(BLUE)🌐 后端API: http://localhost:$(BACKEND_PORT)$(RESET)"
	@echo "$(BLUE)📚 API文档: http://localhost:$(BACKEND_PORT)/api/v1$(RESET)"
	@echo "$(YELLOW)💡 使用 'make stop' 停止服务$(RESET)"

## 快速HTTPS启动
quick-https: 
	@echo "$(GREEN)🚀 快速HTTPS启动...$(RESET)"
	@echo "$(YELLOW)选择启动模式:$(RESET)"
	@echo "1) HTTPS模式 (推荐)"
	@echo "2) HTTP模式"
	@read -p "请选择 (1-2): " choice; \
	case $$choice in \
		1) make up-https ;; \
		2) make up ;; \
		*) echo "$(RED)无效选择$(RESET)" ;; \
	esac

## 配置Let's Encrypt证书
setup-letsencrypt:
	@echo "$(GREEN)🔐 配置Let's Encrypt SSL证书...$(RESET)"
	@echo "$(YELLOW)请输入您的域名 (默认: gogotou.cn):$(RESET)"
	@read -p "域名: " domain; \
	if [ -z "$$domain" ]; then domain="gogotou.cn"; fi; \
	echo "$(YELLOW)请输入您的邮箱:$(RESET)"; \
	read -p "邮箱: " email; \
	if [ -z "$$email" ]; then email="admin@$$domain"; fi; \
	./scripts/setup-letsencrypt.sh "$$domain" "$$email"
	@echo "$(GREEN)✅ Let's Encrypt证书配置完成$(RESET)"

## 构建生产环境镜像
build-prod:
	@echo "$(GREEN)🔨 构建生产环境Docker镜像...$(RESET)"
	docker build -t $(FRONTEND_IMAGE):$(VERSION)-prod -t $(FRONTEND_IMAGE):latest-prod -f $(FRONTEND_PROD_DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 生产环境镜像构建完成: $(FRONTEND_IMAGE):$(VERSION)-prod$(RESET)"

## 部署生产环境服务
up-prod: build-backend build-prod
	@echo "$(GREEN)🚀 部署生产环境HTTPS服务...$(RESET)"
	@echo "$(YELLOW)正在创建网络...$(RESET)"
	docker network create zhitou-network 2>/dev/null || true
	@echo "$(YELLOW)正在启动Go后端服务...$(RESET)"
	docker run -d \
		--name $(BACKEND_CONTAINER) \
		--network zhitou-network \
		-p $(BACKEND_PORT):$(BACKEND_PORT) \
		-e ENVIRONMENT=production \
		--restart unless-stopped \
		$(BACKEND_IMAGE):latest
	@echo "$(YELLOW)正在启动生产环境前端服务...$(RESET)"
	docker run -d \
		--name $(FRONTEND_CONTAINER)-prod \
		--network zhitou-network \
		-p 80:80 \
		-p 443:443 \
		-v "/etc/letsencrypt:/etc/letsencrypt:ro" \
		--restart unless-stopped \
		$(FRONTEND_IMAGE):latest-prod
	@echo "$(GREEN)✅ 生产环境HTTPS服务部署完成$(RESET)"
	@echo "$(BLUE)🌐 HTTPS访问: https://您的域名$(RESET)"
	@echo "$(BLUE)🌐 HTTP重定向: http://您的域名 -> https://您的域名$(RESET)"
	@echo "$(BLUE)🌐 后端API: http://您的域名:$(BACKEND_PORT)$(RESET)"
	@echo "$(BLUE)📚 API文档: http://您的域名:$(BACKEND_PORT)/docs$(RESET)"
	@echo "$(YELLOW)💡 使用 'make stop' 停止服务$(RESET)"

## 健康检查
health:
	@echo "$(BLUE)🔍 健康检查...$(RESET)"
	@curl -f http://localhost:$(BACKEND_PORT)/health 2>/dev/null && \
		echo "$(GREEN)✅ Go后端服务正常运行$(RESET)" || \
		echo "$(RED)❌ Go后端服务未响应$(RESET)"
	@curl -f http://localhost:80 2>/dev/null && \
		echo "$(GREEN)✅ 前端服务正常运行$(RESET)" || \
		echo "$(RED)❌ 前端服务未响应$(RESET)"

## 运行测试
test: build-backend
	@echo "$(BLUE)🧪 运行Go后端测试...$(RESET)"
	cd backend-go && go test ./... || echo "$(YELLOW)⚠️ 测试失败$(RESET)"

## 使用docker-compose部署
deploy:
	@echo "$(GREEN)🚀 使用docker-compose部署...$(RESET)"
	docker-compose down
	docker-compose build
	docker-compose up -d
	@echo "$(GREEN)✅ 部署完成$(RESET)"
	@echo "$(BLUE)🌐 访问地址: http://localhost:$(BACKEND_PORT)$(RESET)"

## 部署到Kubernetes
k8s:
	@echo "$(GREEN)☸️ 部署到Kubernetes...$(RESET)"
	kubectl apply -f deploy/k8s/k8s-deployment.yaml
	@echo "$(GREEN)✅ K8s部署完成$(RESET)"

## 从Kubernetes删除
k8s-delete:
	@echo "$(YELLOW)🗑 从Kubernetes删除...$(RESET)"
	kubectl delete -f deploy/k8s/k8s-deployment.yaml
	@echo "$(GREEN)✅ K8s删除完成$(RESET)"

## 显示版本信息
version:
	@echo "$(BLUE)版本信息:$(RESET)"
	@echo "  项目名称: $(PROJECT_NAME)"
	@echo "  当前版本: $(VERSION)"
	@echo "  后端镜像: $(BACKEND_IMAGE):$(VERSION)"
	@echo "  前端镜像: $(FRONTEND_IMAGE):$(VERSION)"
	@echo "  Git提交:  $(shell git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "  构建时间: $(shell date)"

## 构建多架构镜像
build-multi:
	@echo "$(GREEN)🏗 构建多架构镜像...$(RESET)"
	docker buildx create --use --name multiarch-builder 2>/dev/null || true
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(BACKEND_IMAGE):$(VERSION) \
		-t $(BACKEND_IMAGE):latest \
		--push \
		-f $(BACKEND_DOCKERFILE) $(DOCKER_CONTEXT)
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(FRONTEND_IMAGE):$(VERSION) \
		-t $(FRONTEND_IMAGE):latest \
		--push \
		-f $(FRONTEND_DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 多架构镜像构建完成$(RESET)"

## 生成发布包
release: build
	@echo "$(GREEN)📦 生成发布包...$(RESET)"
	mkdir -p release
	docker save $(BACKEND_IMAGE):$(VERSION) | gzip > release/$(PROJECT_NAME)-backend-$(VERSION).tar.gz
	docker save $(FRONTEND_IMAGE):$(VERSION) | gzip > release/$(PROJECT_NAME)-frontend-$(VERSION).tar.gz
	cp docker-compose.yml release/ 2>/dev/null || true
	cp deploy/k8s/k8s-deployment.yaml release/ 2>/dev/null || true
	cp README.md release/
	@echo "$(GREEN)✅ 发布包生成完成: release/$(PROJECT_NAME)-*-$(VERSION).tar.gz$(RESET)"

## 加载发布包
load-release:
	@echo "$(GREEN)📥 加载发布包...$(RESET)"
	@if [ -f "release/$(PROJECT_NAME)-backend-$(VERSION).tar.gz" ]; then \
		docker load < release/$(PROJECT_NAME)-backend-$(VERSION).tar.gz; \
		echo "$(GREEN)✅ 后端发布包加载完成$(RESET)"; \
	else \
		echo "$(RED)❌ 后端发布包不存在$(RESET)"; \
	fi
	@if [ -f "release/$(PROJECT_NAME)-frontend-$(VERSION).tar.gz" ]; then \
		docker load < release/$(PROJECT_NAME)-frontend-$(VERSION).tar.gz; \
		echo "$(GREEN)✅ 前端发布包加载完成$(RESET)"; \
	else \
		echo "$(RED)❌ 前端发布包不存在$(RESET)"; \
	fi

## 显示资源使用情况
stats:
	@echo "$(BLUE)📈 资源使用情况:$(RESET)"
	docker stats $(BACKEND_CONTAINER) $(FRONTEND_CONTAINER) --no-stream 2>/dev/null || \
	docker stats $(PROJECT_NAME) --no-stream 2>/dev/null || \
	echo "$(YELLOW)⚠️ 容器未运行$(RESET)"

## 备份数据
backup:
	@echo "$(GREEN)💾 备份数据...$(RESET)"
	mkdir -p backups
	docker exec $(BACKEND_CONTAINER) tar czf - /app/data 2>/dev/null > backups/backend-backup-$(shell date +%Y%m%d-%H%M%S).tar.gz || \
	echo "$(YELLOW)⚠️ 无数据需要备份$(RESET)"

## 显示镜像信息
info:
	@echo "$(BLUE)🔍 镜像信息:$(RESET)"
	@echo "$(YELLOW)后端镜像:$(RESET)"
	docker images $(BACKEND_IMAGE) --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}\t{{.CreatedAt}}"
	@echo "$(YELLOW)前端镜像:$(RESET)"
	docker images $(FRONTEND_IMAGE) --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}\t{{.CreatedAt}}"

## 查看MySQL日志
logs-db:
	@echo "$(BLUE)📝 查看MySQL容器日志...$(RESET)"
	docker logs -f zhitou-mysql

## 进入MySQL容器
db-shell:
	@echo "$(BLUE)🐚 进入MySQL容器...$(RESET)"
	docker exec -it zhitou-mysql mysql -u root -p123456 stock_prediction

## Docker Compose相关命令
compose-up:
	@echo "$(GREEN)🚀 启动Docker Compose服务...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)✅ 服务已启动$(RESET)"

compose-down:
	@echo "$(YELLOW)🛑 停止Docker Compose服务...$(RESET)"
	docker-compose down
	@echo "$(GREEN)✅ 服务已停止$(RESET)"

compose-logs:
	@echo "$(BLUE)📝 查看Docker Compose日志...$(RESET)"
	docker-compose logs -f

compose-ps:
	@echo "$(BLUE)📊 Docker Compose服务状态:$(RESET)"
	docker-compose ps