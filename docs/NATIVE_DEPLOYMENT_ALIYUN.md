# 生产环境部署记录：阿里云原生部署（非容器化）方案

> 本文档记录了在阿里云 ECS（8.155.167.87，域名 gogotou.cn）上的一次实际部署过程。
> 与仓库默认的 Docker Compose 方案不同，本次因服务器容器环境限制，采用了**原生环境部署**（Go 直接编译运行 + Nginx 反代 + MariaDB），供后续维护和类似环境部署参考。

## 适用场景

当目标服务器满足以下情况时，可参考本方案：
- 容器引擎为 **Podman** 而非 Docker（如阿里云 Alibaba Cloud Linux 默认环境）
- 从 Docker Hub 拉取基础镜像（`golang`、`node` 等）网络不通或速度极慢，且没有可用的镜像加速源
- 服务器已有 `yum`/`dnf` 包管理器，可直接安装 Go、Node.js、MariaDB、Nginx 等原生运行时

## 环境信息

| 项目 | 说明 |
|---|---|
| 服务器 | 阿里云 ECS，Alibaba Cloud Linux 3 (OpenAnolis)，公网 IP `8.155.167.87` |
| 域名 | `gogotou.cn`（已解析到上述 IP） |
| 容器引擎 | Podman 4.9（非 Docker，需要 `podman-compose` 才能识别 `docker-compose.yml`） |
| 部署方式 | 原生部署（Go 二进制 + Nginx 静态资源/反代 + MariaDB） |

## 部署步骤

### 1. 排查容器方案不可行

```bash
docker --version && docker compose version
# 实际返回 podman 4.9.4-rhel，且缺少 compose 插件

yum install -y podman-compose
podman-compose build
# golang / node 基础镜像从 docker.io 拉取全部超时（阿里云内网访问 Docker Hub 及国内加速镜像站均失败）
```

结论：直接放弃容器化构建，改为在服务器上安装原生运行时，编译/构建后直接运行。

### 2. 安装原生运行时

```bash
yum install -y golang mariadb-server nginx certbot python3-certbot-nginx
```

服务器自带 Node.js（如没有需额外安装 `nodejs`）。

### 3. 初始化数据库（MariaDB）

```bash
systemctl start mariadb && systemctl enable mariadb

mysql -u root << 'EOF'
CREATE DATABASE IF NOT EXISTS stock_prediction CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';  -- 默认走 socket 认证，需显式设密码才能支持 TCP 连接
FLUSH PRIVILEGES;
EOF
```

> ⚠️ 生产环境务必把密码改成强密码，这里为了跑通流程使用了与 `docker-compose.yml` 一致的弱密码 `123456`。

### 4. 编译并启动后端（Go）

```bash
# Go 官方 proxy.golang.org 在国内被墙，需切换国内镜像
go env -w GOPROXY=https://goproxy.cn,direct

cd backend-go
go mod tidy
go build -o stock-backend cmd/main.go

export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=stock_prediction
export PORT=8000
export ENVIRONMENT=production
nohup ./stock-backend > /root/backend.log 2>&1 &
```

后端会自动完成建表（GORM AutoMigrate）并启动每日 15:10 的预测定时任务。

> ⚠️ **待改进**：目前用 `nohup` 挂后台运行，服务器重启后进程不会自动恢复。生产环境建议配置为 `systemd` service（`/etc/systemd/system/stock-backend.service`），实现进程守护和开机自启。

### 5. 构建前端（Vite）并部署静态资源

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
npm install
npm run build

# 注意：产物不要放在 /root 下，Nginx worker 进程权限不足会读取失败（403）
mkdir -p /var/www/stock_prediction
cp -r dist/* /var/www/stock_prediction/
chmod -R 755 /var/www/stock_prediction
```

### 6. 配置 Nginx（静态资源 + API 反代）

创建 `/etc/nginx/conf.d/stock_prediction.conf`：

```nginx
server {
    listen 80;
    server_name gogotou.cn;
    root /var/www/stock_prediction;
    index index.html;

    location ~* \.html$ {
        expires -1;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }

    location ~* \.(js|css)$ {
        expires 1h;
        add_header Cache-Control "public, must-revalidate";
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8000/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

> ⚠️ **重要坑点**：Alibaba Cloud Linux 的 `/etc/nginx/nginx.conf` 主配置文件自带一个监听 80 端口、`server_name _` 的默认 server 块（约第 37-53 行），会与 `conf.d/` 下的自定义配置冲突（`conflicting server name` 警告，且默认块优先生效导致 403）。需要注释掉该默认 server 块：
>
> ```bash
> sed -i '37,53s/^/#/' /etc/nginx/nginx.conf
> ```
> 请先用 `grep -n "server {" -A 20 /etc/nginx/nginx.conf` 确认实际行号，不同系统版本可能不同。

```bash
rm -f /etc/nginx/conf.d/default.conf
nginx -t
systemctl start nginx && systemctl enable nginx
```

### 7. 配置 HTTPS（Let's Encrypt，域名方式）

Let's Encrypt 证书**必须绑定域名**，无法直接给纯 IP 签发。前提：域名已通过 DNS A 记录解析到服务器公网 IP。

```bash
# 确认域名已解析
dig +short gogotou.cn A   # 应返回服务器公网 IP

# 用 certbot 的 nginx 插件自动申请并改写 nginx 配置
certbot --nginx -d gogotou.cn --email admin@gogotou.cn --agree-tos --non-interactive
```

certbot 会自动：
- 在 443 端口新增 SSL server 块，写入证书路径
- 把 80 端口的 HTTP 请求 301 重定向到 HTTPS
- 注册 systemd timer / cron 任务，定期自动续期（证书有效期 90 天）

> 若还需要支持 `www.gogotou.cn`，需先在 DNS 加一条 A 记录指向服务器 IP，再重新执行：
> `certbot --nginx -d gogotou.cn -d www.gogotou.cn --email admin@gogotou.cn --agree-tos --non-interactive`

### 8. 验证

```bash
curl -s -o /dev/null -w "%{http_code}\n" https://gogotou.cn/          # 200
curl -s https://gogotou.cn/api/v1/prediction-stats                     # 返回 JSON
echo | openssl s_client -connect gogotou.cn:443 -servername gogotou.cn 2>/dev/null | openssl x509 -noout -dates
```

## 最终结果

- 访问地址：**https://gogotou.cn/**（HTTP 自动跳转 HTTPS）
- 后端 API：`https://gogotou.cn/api/v1/...`（经 Nginx 反代到本机 `127.0.0.1:8000`）
- 证书：Let's Encrypt，自动续期已配置

## 已知待办 / 后续改进

1. **后端进程守护**：当前用 `nohup` 后台运行，重启服务器后不会自动恢复，需要迁移到 `systemd` service。
2. **数据库密码加固**：当前使用弱密码 `123456`（沿用仓库 `docker-compose.yml` 默认值），生产环境应更换为强密码并同步更新后端环境变量。
3. **www 子域名**：`www.gogotou.cn` 尚未配置 DNS 和证书，如有需要请按上文补充。
4. **与仓库既有 Docker 方案的关系**：本文档描述的是绕过 Docker/Podman 的应急方案，仅在容器构建不可行时使用；如目标服务器容器网络通畅，仍建议优先使用 `DEPLOYMENT.md` 中的 `make deploy` 方案。
