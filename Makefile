# gogotou - AI股市指数预测平台 Makefile
# 后端: Python/FastAPI   前端: Vue 3

PROJECT_NAME := zhitou-prediction
BACKEND_IMAGE := zhitou-prediction-backend
FRONTEND_IMAGE := zhitou-prediction-frontend

# 端口
BACKEND_PORT := 8001
FRONTEND_PORT := 9000

DOCKER_COMPOSE_V2 := docker compose
DOCKER_COMPOSE_V1 := docker-compose

GREEN  := \033[32m
YELLOW := \033[33m
BLUE   := \033[34m
RESET  := \033[0m

.PHONY: help install-deps backend frontend dev build build-backend build-frontend up down logs health test clean

.DEFAULT_GOAL := help

help:
	@echo "$(BLUE)gogotou - AI指数预测平台 (Python后端 + Vue前端)$(RESET)"
	@echo "==============================================="
	@echo "$(GREEN)开发:$(RESET)"
	@echo "  install-deps - 安装前后端依赖"
	@echo "  backend      - 启动Python后端(8001, 热重载)"
	@echo "  frontend     - 启动Vue前端(9000)"
	@echo "  dev          - 同时启动前后端"
	@echo "$(GREEN)Docker:$(RESET)"
	@echo "  build        - 构建前后端镜像"
	@echo "  up           - 启动全部服务 (docker compose)"
	@echo "  down         - 停止全部服务"
	@echo "  logs         - 查看日志"
	@echo "$(GREEN)其他:$(RESET)"
	@echo "  health       - 后端健康检查"
	@echo "  test         - 后端冒烟测试"
	@echo "  clean        - 清理构建产物"

install-deps:
	@echo "$(GREEN)📦 安装后端依赖...$(RESET)"
	cd backend && pip install -r requirements.txt
	@echo "$(GREEN)📦 安装前端依赖...$(RESET)"
	cd frontend && npm install

backend:
	@echo "$(GREEN)🚀 启动Python后端 (http://localhost:$(BACKEND_PORT))$(RESET)"
	cd backend && python3 -m uvicorn app.main:app --host 0.0.0.0 --port $(BACKEND_PORT) --reload

frontend:
	@echo "$(GREEN)🚀 启动Vue前端 (http://localhost:$(FRONTEND_PORT))$(RESET)"
	cd frontend && npm run dev

dev:
	@echo "$(GREEN)🚀 启动前后端开发环境...$(RESET)"
	@echo "$(YELLOW)后端: http://localhost:$(BACKEND_PORT)  前端: http://localhost:$(FRONTEND_PORT)$(RESET)"
	make -j2 backend frontend

build:
	@echo "$(GREEN)🔨 构建镜像...$(RESET)"
	$(DOCKER_COMPOSE_V2) build

build-backend:
	docker build -t $(BACKEND_IMAGE):latest -f deployment/backend/Dockerfile .

build-frontend:
	docker build -t $(FRONTEND_IMAGE):latest -f deployment/frontend/Dockerfile .

up:
	@echo "$(GREEN)🚀 启动服务 (docker compose)...$(RESET)"
	$(DOCKER_COMPOSE_V2) up -d --build
	@echo "$(BLUE)前端: http://localhost:80  后端: http://localhost:$(BACKEND_PORT)$(RESET)"

down:
	$(DOCKER_COMPOSE_V2) down

logs:
	$(DOCKER_COMPOSE_V2) logs -f

health:
	@curl -f http://localhost:$(BACKEND_PORT)/health >/dev/null 2>&1 && \
		echo "$(GREEN)✅ 后端正常$(RESET)" || echo "$(YELLOW)❌ 后端未响应$(RESET)"

test:
	@echo "$(BLUE)🧪 后端冒烟测试...$(RESET)"
	cd backend && python3 -c "import app.main; print('import ok')"

clean:
	rm -rf frontend/dist backend/__pycache__ backend/*.pyc
