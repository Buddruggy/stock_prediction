# 智投预测 V2.0 - 前端 Docker镜像
# 多阶段构建：构建阶段 + 运行阶段

# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制package文件
COPY frontend/package*.json ./

# 安装依赖 - 使用国内镜像源加速
RUN npm config set registry https://registry.npmmirror.com && \
    npm install

# 复制前端源代码
COPY frontend/ .

# 构建生产版本
RUN npm run build

# 运行阶段
FROM nginx:alpine

# 复制自定义nginx配置
COPY deployment/frontend/nginx.conf /etc/nginx/nginx.conf

# 复制构建好的文件到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 静态资源已经包含在Vue构建产物中

# 创建非root用户
RUN addgroup -g 1001 -S app && \
    adduser -S app -u 1001 -G app && \
    chown -R app:app /usr/share/nginx/html && \
    chown -R app:app /var/cache/nginx && \
    chown -R app:app /var/log/nginx && \
    chown -R app:app /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R app:app /var/run/nginx.pid

USER app

# 健康检查
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:9000/ || exit 1

# 暴露端口
EXPOSE 9000

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
