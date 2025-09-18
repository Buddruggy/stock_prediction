# gogotou - 部署指南

## 环境说明

本项目已针对不同环境的Docker Compose版本进行了优化配置，请根据您的环境选择合适的部署命令。

## 部署命令

### 🖥️ 服务器环境 (推荐)
适用于：Linux服务器、生产环境、新版本Docker环境

**使用命令：**
```bash
make deploy
```

**技术说明：**
- 使用 `docker compose` (Docker Compose v2)
- 适合生产环境部署
- 性能更好，启动更快

### 💻 Mac本地环境
适用于：Mac开发环境、旧版本Docker环境

**使用命令：**
```bash
make deploy-mac
```

**技术说明：**
- 使用 `docker-compose` (docker-compose v1)
- 适合本地开发和测试
- 兼容性更好

## 管理命令

### 服务器环境管理命令
```bash
# 启动服务
make compose-up

# 停止服务  
make compose-down

# 查看日志
make compose-logs

# 查看状态
make compose-ps
```

### Mac环境管理命令
```bash
# 启动服务
make compose-up-mac

# 停止服务
make compose-down-mac

# 查看日志
make compose-logs-mac

# 查看状态
make compose-ps-mac
```

## 完整服务架构

部署成功后，您将获得以下服务：

1. **MySQL 8.0 数据库** (端口 3306)
   - 数据库名：`stock_prediction`
   - 用户名：`root`，密码：`123456`
   - 统一表结构：`predictions` 和 `historical_data`

2. **Go 后端服务** (端口 8000)
   - API 端点：http://localhost:8000/api/v1
   - 健康检查：http://localhost:8000/health

3. **前端服务** (端口 80 和 9000)
   - 主要访问：http://localhost:80
   - 备用访问：http://localhost:9000

## 数据持久化功能

- ✅ **预测数据存储**：每日预测结果自动存储到MySQL
- ✅ **历史数据存储**：股票历史数据持久化
- ✅ **统一表结构**：使用 `index_code` 字段区分不同指数
- ✅ **定时任务**：下午3:10分（A股收盘后）自动刷新数据

## 数据优先级策略

系统实现智能数据获取策略：
**数据库 → 缓存 → 实时获取**

前端获取数据时：
1. 优先从数据库查询历史数据
2. 如果数据库无数据或过期，从缓存获取
3. 最后进行实时数据获取

## 故障排除

### 如果部署失败
1. 检查Docker是否正常运行
2. 确保端口 3306、8000、80、9000 未被占用
3. 清理之前的容器：`make stop`
4. 重新部署

### 查看详细日志
```bash
# 服务器环境
make compose-logs

# Mac环境  
make compose-logs-mac

# 查看MySQL日志
make logs-db
```

### 进入数据库
```bash
# 进入MySQL容器
make db-shell
```

## 版本兼容性

| 环境 | Docker Compose版本 | 部署命令 |
|------|-------------------|----------|
| 🖥️ 服务器/生产环境 | v2 (`docker compose`) | `make deploy` |
| 💻 Mac本地环境 | v1 (`docker-compose`) | `make deploy-mac` |

选择正确的命令可以避免版本兼容性问题，确保稳定部署！

---

**提示**：如果不确定您的环境使用哪个版本，可以运行以下命令检查：

```bash
# 检查是否支持新版本
docker compose version

# 检查是否支持旧版本  
docker-compose version
```