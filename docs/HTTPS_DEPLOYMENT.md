# HTTPS 部署指南

## 🔐 解决Chrome"不安全"警告

Chrome浏览器显示"不安全"是因为使用了HTTP协议。本指南提供两种HTTPS解决方案：

### 方案1: 自签名证书（测试环境）

#### 1. 生成SSL证书
```bash
# 生成自签名证书
make ssl-cert

# 或手动生成
./scripts/generate-ssl-cert.sh yourdomain.com
```

#### 2. 部署HTTPS服务
```bash
# 一键部署HTTPS服务
make up-https

# 或分步执行
make build-frontend-https
make up-https
```

#### 3. 访问测试
- HTTPS: `https://localhost:443`
- HTTP自动重定向到HTTPS: `http://localhost:80`

⚠️ **注意**: 自签名证书浏览器会显示安全警告，需要手动信任。

### 方案2: Let's Encrypt免费证书（生产环境）

#### 1. 安装Certbot
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install certbot

# CentOS/RHEL
sudo yum install certbot

# macOS
brew install certbot
```

#### 2. 获取证书
```bash
# 停止当前服务
make stop

# 获取Let's Encrypt证书
sudo certbot certonly --standalone -d yourdomain.com

# 证书位置: /etc/letsencrypt/live/yourdomain.com/
```

#### 3. 配置生产环境nginx
创建生产环境nginx配置：

```nginx
# deployment/frontend/nginx-production.conf
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    # HTTP重定向到HTTPS
    server {
        listen 80;
        server_name yourdomain.com www.yourdomain.com;
        return 301 https://$host$request_uri;
    }
    
    # HTTPS服务器
    server {
        listen 443 ssl http2;
        server_name yourdomain.com www.yourdomain.com;
        
        # Let's Encrypt证书
        ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
        
        # SSL安全配置
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers on;
        
        # 安全头
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options DENY always;
        add_header X-Content-Type-Options nosniff always;
        
        # 其他配置...
        root /usr/share/nginx/html;
        index index.html;
        
        location / {
            try_files $uri $uri/ /index.html;
        }
        
        location /api/ {
            proxy_pass http://zhitou-prediction-backend:8000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

#### 4. 生产环境部署
```bash
# 构建生产镜像
docker build -t zhitou-frontend:prod -f deployment/frontend/Dockerfile.prod .

# 运行生产容器
docker run -d \
  --name zhitou-frontend-prod \
  --network zhitou-network \
  -p 80:80 \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  --restart unless-stopped \
  zhitou-frontend:prod
```

#### 5. 自动续期
```bash
# 设置自动续期
sudo crontab -e

# 添加以下行（每月1号检查续期）
0 0 1 * * certbot renew --quiet && docker restart zhitou-frontend-prod
```

## 🚀 快速开始

### 测试环境
```bash
# 1. 生成自签名证书
make ssl-cert

# 2. 部署HTTPS服务
make up-https

# 3. 访问 https://localhost:443
```

### 生产环境
```bash
# 1. 获取Let's Encrypt证书
sudo certbot certonly --standalone -d yourdomain.com

# 2. 使用生产配置部署
# (需要创建生产环境Dockerfile和nginx配置)

# 3. 访问 https://yourdomain.com
```

## 🔧 故障排除

### 1. 证书问题
```bash
# 检查证书有效性
openssl x509 -in ssl/cert.pem -text -noout

# 检查证书过期时间
openssl x509 -in ssl/cert.pem -dates -noout
```

### 2. 端口冲突
```bash
# 检查端口占用
sudo netstat -tlnp | grep :443
sudo netstat -tlnp | grep :80

# 停止占用端口的服务
sudo systemctl stop apache2  # 如果Apache占用80端口
sudo systemctl stop nginx     # 如果Nginx占用80端口
```

### 3. 防火墙设置
```bash
# Ubuntu/Debian
sudo ufw allow 80
sudo ufw allow 443

# CentOS/RHEL
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

## 📋 检查清单

- [ ] 域名已解析到服务器IP
- [ ] 80和443端口已开放
- [ ] SSL证书已正确配置
- [ ] HTTP自动重定向到HTTPS
- [ ] 安全头已配置
- [ ] 证书自动续期已设置

## 🎯 最佳实践

1. **使用HTTPS**: 所有生产环境都应该使用HTTPS
2. **自动重定向**: HTTP请求自动重定向到HTTPS
3. **安全头**: 配置适当的安全头
4. **证书管理**: 设置自动续期
5. **监控**: 监控证书过期时间

完成HTTPS配置后，Chrome将不再显示"不安全"警告！🎉
