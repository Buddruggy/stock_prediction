#!/bin/bash

# 生成自签名SSL证书脚本
# 用于测试环境，生产环境建议使用Let's Encrypt

set -e

CERT_DIR="ssl"
DOMAIN="${1:-localhost}"

echo "🔐 生成SSL证书..."
echo "域名: $DOMAIN"

# 创建证书目录
mkdir -p "$CERT_DIR"

# 生成私钥
openssl genrsa -out "$CERT_DIR/key.pem" 2048

# 生成证书签名请求
openssl req -new -key "$CERT_DIR/key.pem" -out "$CERT_DIR/cert.csr" -subj "/C=CN/ST=Beijing/L=Beijing/O=Zhitou/OU=IT/CN=$DOMAIN"

# 生成自签名证书
openssl x509 -req -days 365 -in "$CERT_DIR/cert.csr" -signkey "$CERT_DIR/key.pem" -out "$CERT_DIR/cert.pem"

# 清理临时文件
rm "$CERT_DIR/cert.csr"

echo "✅ SSL证书生成完成！"
echo "📁 证书位置: $CERT_DIR/"
echo "   - cert.pem: 证书文件"
echo "   - key.pem: 私钥文件"
echo ""
echo "⚠️  注意: 这是自签名证书，浏览器会显示安全警告"
echo "   生产环境建议使用 Let's Encrypt 免费证书"
