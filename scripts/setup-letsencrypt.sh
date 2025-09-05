#!/bin/bash

# Let's Encrypt SSL证书配置脚本
# 用于生产环境获取免费SSL证书

set -e

DOMAIN="${1:-gogotou.cn}"
EMAIL="${2:-admin@gogotou.cn}"

echo "🔐 配置Let's Encrypt SSL证书"
echo "域名: $DOMAIN"
echo "邮箱: $EMAIL"
echo ""

# 检查是否安装了certbot
if ! command -v certbot &> /dev/null; then
    echo "❌ 未安装certbot，正在安装..."
    
    # 检测操作系统
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux系统
        if command -v apt &> /dev/null; then
            # Ubuntu/Debian
            sudo apt update
            sudo apt install -y certbot
        elif command -v yum &> /dev/null; then
            # CentOS/RHEL
            sudo yum install -y certbot
        else
            echo "❌ 不支持的Linux发行版"
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install certbot
        else
            echo "❌ 请先安装Homebrew: https://brew.sh"
            exit 1
        fi
    else
        echo "❌ 不支持的操作系统"
        exit 1
    fi
fi

echo "✅ certbot已安装"

# 停止当前服务（释放80端口）
echo "🛑 停止当前服务..."
make stop 2>/dev/null || true

# 获取证书
echo "🔐 获取Let's Encrypt证书..."
sudo certbot certonly --standalone -d "$DOMAIN" --email "$EMAIL" --agree-tos --non-interactive

# 检查证书是否获取成功
if [ -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
    echo "✅ SSL证书获取成功！"
    echo "📁 证书位置: /etc/letsencrypt/live/$DOMAIN/"
    echo "   - fullchain.pem: 完整证书链"
    echo "   - privkey.pem: 私钥"
    
    # 创建生产环境nginx配置
    echo "📝 创建生产环境nginx配置..."
    cat > deployment/frontend/nginx-production.conf << EOF
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # 日志格式
    log_format main '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                    '\$status \$body_bytes_sent "\$http_referer" '
                    '"\$http_user_agent" "\$http_x_forwarded_for"';
    
    access_log /var/log/nginx/access.log main;
    
    # 基本设置
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    
    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript 
               application/javascript application/xml+rss 
               application/json application/xml;
    
    # HTTP重定向到HTTPS
    server {
        listen 80;
        server_name $DOMAIN www.$DOMAIN;
        return 301 https://\$host\$request_uri;
    }
    
    # HTTPS服务器配置
    server {
        listen 443 ssl http2;
        server_name $DOMAIN www.$DOMAIN;
        root /usr/share/nginx/html;
        index index.html;
        
        # Let's Encrypt证书
        ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
        
        # SSL安全配置
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384;
        ssl_prefer_server_ciphers on;
        ssl_session_cache shared:SSL:10m;
        ssl_session_timeout 10m;
        
        # 安全头
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options DENY always;
        add_header X-Content-Type-Options nosniff always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;
        
        # 静态资源缓存
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
        
        # API代理到后端
        location /api/ {
            proxy_pass http://zhitou-prediction-backend:8000;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }
        
        # SPA路由处理
        location / {
            try_files \$uri \$uri/ /index.html;
        }
        
        # 健康检查
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
EOF
    
    echo "✅ 生产环境nginx配置已创建"
    echo ""
    echo "🚀 下一步操作:"
    echo "1. 使用生产配置构建镜像:"
    echo "   docker build -t zhitou-frontend:prod -f deployment/frontend/Dockerfile.prod ."
    echo ""
    echo "2. 运行生产容器:"
    echo "   docker run -d --name zhitou-frontend-prod --network zhitou-network -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro --restart unless-stopped zhitou-frontend:prod"
    echo ""
    echo "3. 设置自动续期:"
    echo "   sudo crontab -e"
    echo "   # 添加: 0 0 1 * * certbot renew --quiet && docker restart zhitou-frontend-prod"
    
else
    echo "❌ SSL证书获取失败"
    echo "请检查:"
    echo "1. 域名是否正确解析到服务器IP"
    echo "2. 80端口是否开放"
    echo "3. 是否有其他服务占用80端口"
    exit 1
fi
