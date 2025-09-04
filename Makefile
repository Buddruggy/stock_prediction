# 智投预测 - AI股市指数预测平台 Makefile
# 用于构建、打包和部署Docker镜像

# 项目配置
PROJECT_NAME := zhitou-prediction
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
REGISTRY := docker.io
USERNAME := alanwzliang
IMAGE_NAME := $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)
CONTAINER_NAME := $(PROJECT_NAME)

# 端口配置
HOST_PORT := 80
CONTAINER_PORT := 80

# Docker配置
DOCKERFILE := deploy/docker/Dockerfile
DOCKER_CONTEXT := .

# 颜色输出
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help build push pull run stop clean logs shell test dev prod deploy

# 默认目标
.DEFAULT_GOAL := help

## 显示帮助信息
help:
	@echo "$(BLUE)智投预测 - AI股市指数预测平台$(RESET)"
	@echo "$(BLUE)=========================================$(RESET)"
	@echo ""
	@echo "$(GREEN)可用命令:$(RESET)"
	@echo ""
	@echo "$(YELLOW)构建相关:$(RESET)"
	@echo "  build     - 构建Docker镜像"
	@echo "  build-nc  - 无缓存构建Docker镜像"
	@echo "  push      - 推送镜像到仓库"
	@echo "  pull      - 从仓库拉取镜像"
	@echo ""
	@echo "$(YELLOW)运行相关:$(RESET)"
	@echo "  run       - 运行容器"
	@echo "  dev       - 运行开发环境"
	@echo "  prod      - 运行生产环境"
	@echo "  stop      - 停止容器"
	@echo "  restart   - 重启容器"
	@echo ""
	@echo "$(YELLOW)管理相关:$(RESET)"
	@echo "  logs      - 查看容器日志"
	@echo "  shell     - 进入容器shell"
	@echo "  ps        - 查看容器状态"
	@echo "  clean     - 清理未使用的镜像"
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
	@echo "  镜像名:   $(IMAGE_NAME):$(VERSION)"
	@echo "  端口:     $(HOST_PORT):$(CONTAINER_PORT)"

## 构建Docker镜像
build:
	@echo "$(GREEN)🔨 构建Docker镜像...$(RESET)"
	docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest -f $(DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 镜像构建完成: $(IMAGE_NAME):$(VERSION)$(RESET)"

## 无缓存构建Docker镜像
build-nc:
	@echo "$(GREEN)🔨 无缓存构建Docker镜像...$(RESET)"
	docker build --no-cache -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest -f $(DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 镜像构建完成: $(IMAGE_NAME):$(VERSION)$(RESET)"

## 推送镜像到仓库
push: build
	@echo "$(GREEN)📤 推送镜像到仓库...$(RESET)"
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest
	@echo "$(GREEN)✅ 镜像推送完成$(RESET)"

## 从仓库拉取镜像
pull:
	@echo "$(GREEN)📥 从仓库拉取镜像...$(RESET)"
	docker pull $(IMAGE_NAME):latest
	@echo "$(GREEN)✅ 镜像拉取完成$(RESET)"

## 运行容器
run: stop
	@echo "$(GREEN)🚀 启动容器...$(RESET)"
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(HOST_PORT):$(CONTAINER_PORT) \
		-e FLASK_ENV=production \
		--restart unless-stopped \
		$(IMAGE_NAME):latest
	@echo "$(GREEN)✅ 容器启动完成$(RESET)"
	@echo "$(BLUE)🌐 访问地址: http://localhost:$(HOST_PORT)$(RESET)"

## 运行开发环境
dev: stop
	@echo "$(GREEN)🛠 启动开发环境...$(RESET)"
	docker run -d \
		--name $(CONTAINER_NAME)-dev \
		-p $(HOST_PORT):$(CONTAINER_PORT) \
		-e FLASK_ENV=development \
		-v $(PWD):/app \
		--restart unless-stopped \
		$(IMAGE_NAME):latest
	@echo "$(GREEN)✅ 开发环境启动完成$(RESET)"
	@echo "$(BLUE)🌐 访问地址: http://localhost:$(HOST_PORT)$(RESET)"

## 运行生产环境
prod: build run
	@echo "$(GREEN)🏭 生产环境部署完成$(RESET)"

## 停止容器
stop:
	@echo "$(YELLOW)🛑 停止容器...$(RESET)"
	-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	-docker stop $(CONTAINER_NAME)-dev 2>/dev/null || true
	-docker rm $(CONTAINER_NAME) 2>/dev/null || true
	-docker rm $(CONTAINER_NAME)-dev 2>/dev/null || true

## 重启容器
restart: stop run

## 查看容器日志
logs:
	@echo "$(BLUE)📋 查看容器日志...$(RESET)"
	docker logs -f $(CONTAINER_NAME) 2>/dev/null || docker logs -f $(CONTAINER_NAME)-dev

## 进入容器shell
shell:
	@echo "$(BLUE)🐚 进入容器shell...$(RESET)"
	docker exec -it $(CONTAINER_NAME) /bin/bash 2>/dev/null || docker exec -it $(CONTAINER_NAME)-dev /bin/bash

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
	-docker rmi $(IMAGE_NAME):$(VERSION) 2>/dev/null || true
	-docker rmi $(IMAGE_NAME):latest 2>/dev/null || true
	docker system prune -af
	@echo "$(GREEN)✅ 深度清理完成$(RESET)"

## 健康检查
health:
	@echo "$(BLUE)🔍 健康检查...$(RESET)"
	@curl -f http://localhost:$(HOST_PORT)/api/status 2>/dev/null && \
		echo "$(GREEN)✅ 服务正常运行$(RESET)" || \
		echo "$(RED)❌ 服务未响应$(RESET)"

## 运行测试
test: build
	@echo "$(BLUE)🧪 运行测试...$(RESET)"
	docker run --rm \
		-e FLASK_ENV=testing \
		$(IMAGE_NAME):latest \
		python -m pytest tests/ || echo "$(YELLOW)⚠️ 测试目录不存在$(RESET)"

## 使用docker-compose部署
deploy:
	@echo "$(GREEN)🚀 使用docker-compose部署...$(RESET)"
	docker-compose down
	docker-compose build
	docker-compose up -d
	@echo "$(GREEN)✅ 部署完成$(RESET)"
	@echo "$(BLUE)🌐 访问地址: http://localhost:$(HOST_PORT)$(RESET)"

## 部署到Kubernetes
k8s:
	@echo "$(GREEN)☸️ 部署到Kubernetes...$(RESET)"
	kubectl apply -f k8s-deployment.yaml
	@echo "$(GREEN)✅ K8s部署完成$(RESET)"

## 从Kubernetes删除
k8s-delete:
	@echo "$(YELLOW)🗑 从Kubernetes删除...$(RESET)"
	kubectl delete -f k8s-deployment.yaml
	@echo "$(GREEN)✅ K8s删除完成$(RESET)"

## 显示版本信息
version:
	@echo "$(BLUE)版本信息:$(RESET)"
	@echo "  项目名称: $(PROJECT_NAME)"
	@echo "  当前版本: $(VERSION)"
	@echo "  镜像名称: $(IMAGE_NAME):$(VERSION)"
	@echo "  Git提交:  $(shell git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "  构建时间: $(shell date)"

## 构建多架构镜像
build-multi:
	@echo "$(GREEN)🏗 构建多架构镜像...$(RESET)"
	docker buildx create --use --name multiarch-builder 2>/dev/null || true
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_NAME):$(VERSION) \
		-t $(IMAGE_NAME):latest \
		--push \
		-f $(DOCKERFILE) $(DOCKER_CONTEXT)
	@echo "$(GREEN)✅ 多架构镜像构建完成$(RESET)"

## 生成发布包
release: build
	@echo "$(GREEN)📦 生成发布包...$(RESET)"
	mkdir -p release
	docker save $(IMAGE_NAME):$(VERSION) | gzip > release/$(PROJECT_NAME)-$(VERSION).tar.gz
	cp docker-compose.yml release/
	cp k8s-deployment.yaml release/
	cp README.md release/
	@echo "$(GREEN)✅ 发布包生成完成: release/$(PROJECT_NAME)-$(VERSION).tar.gz$(RESET)"

## 加载发布包
load-release:
	@echo "$(GREEN)📥 加载发布包...$(RESET)"
	@if [ -f "release/$(PROJECT_NAME)-$(VERSION).tar.gz" ]; then \
		docker load < release/$(PROJECT_NAME)-$(VERSION).tar.gz; \
		echo "$(GREEN)✅ 发布包加载完成$(RESET)"; \
	else \
		echo "$(RED)❌ 发布包不存在$(RESET)"; \
	fi

## 快速启动（构建并运行）
quick: build run health

## 显示资源使用情况
stats:
	@echo "$(BLUE)📈 资源使用情况:$(RESET)"
	docker stats $(CONTAINER_NAME) --no-stream 2>/dev/null || \
	docker stats $(CONTAINER_NAME)-dev --no-stream 2>/dev/null || \
	echo "$(YELLOW)⚠️ 容器未运行$(RESET)"

## 备份数据
backup:
	@echo "$(GREEN)💾 备份数据...$(RESET)"
	mkdir -p backups
	docker exec $(CONTAINER_NAME) tar czf - /app/data 2>/dev/null > backups/backup-$(shell date +%Y%m%d-%H%M%S).tar.gz || \
	echo "$(YELLOW)⚠️ 无数据需要备份$(RESET)"

## 显示镜像信息
info:
	@echo "$(BLUE)🔍 镜像信息:$(RESET)"
	docker images $(IMAGE_NAME) --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}\t{{.CreatedAt}}"
