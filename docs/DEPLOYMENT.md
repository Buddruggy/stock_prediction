# 智投预测 - 云原生部署指南

## 🐳 Docker部署

### 快速开始

```bash
# 1. 使用Docker Compose (推荐)
./docker-run.sh

# 2. 或者手动构建和运行
docker build -t zhitou-prediction .
docker run -p 9000:9000 zhitou-prediction
```

### Docker命令

```bash
# 构建镜像
docker build -t zhitou-prediction:latest .

# 运行容器
docker run -d \
  --name zhitou-prediction \
  -p 9000:9000 \
  -e FLASK_ENV=production \
  zhitou-prediction:latest

# 查看日志
docker logs -f zhitou-prediction

# 停止容器
docker stop zhitou-prediction
```

## ☸️ Kubernetes部署

### 前提条件

- Kubernetes集群 (v1.19+)
- kubectl工具
- Docker镜像已推送到镜像仓库

### 部署步骤

```bash
# 1. 构建并推送镜像到镜像仓库
docker build -t your-registry/zhitou-prediction:v1.0 .
docker push your-registry/zhitou-prediction:v1.0

# 2. 更新k8s-deployment.yaml中的镜像地址
# 将 image: zhitou-prediction:latest 
# 改为 image: your-registry/zhitou-prediction:v1.0

# 3. 部署到Kubernetes
kubectl apply -f k8s-deployment.yaml

# 4. 查看部署状态
kubectl get pods -l app=zhitou-prediction
kubectl get services

# 5. 查看服务日志
kubectl logs -l app=zhitou-prediction -f
```

### 服务访问

```bash
# 获取服务外部IP
kubectl get service zhitou-prediction-service

# 端口转发 (本地测试)
kubectl port-forward service/zhitou-prediction-service 9000:80
```

## ☁️ 云平台部署

### 阿里云容器服务 ACK

```bash
# 1. 登录阿里云容器镜像服务
docker login --username=your-username registry.cn-hangzhou.aliyuncs.com

# 2. 构建并推送镜像
docker build -t registry.cn-hangzhou.aliyuncs.com/your-namespace/zhitou-prediction:v1.0 .
docker push registry.cn-hangzhou.aliyuncs.com/your-namespace/zhitou-prediction:v1.0

# 3. 在ACK控制台创建应用或使用kubectl部署
```

### 腾讯云容器服务 TKE

```bash
# 1. 登录腾讯云容器镜像服务
docker login --username=your-username ccr.ccs.tencentyun.com

# 2. 构建并推送镜像
docker build -t ccr.ccs.tencentyun.com/your-namespace/zhitou-prediction:v1.0 .
docker push ccr.ccs.tencentyun.com/your-namespace/zhitou-prediction:v1.0

# 3. 在TKE控制台部署或使用kubectl
```

### 华为云容器引擎 CCE

```bash
# 1. 登录华为云容器镜像服务
docker login -u cn-north-4@your-iam-user -p your-password swr.cn-north-4.myhuaweicloud.com

# 2. 构建并推送镜像
docker build -t swr.cn-north-4.myhuaweicloud.com/your-namespace/zhitou-prediction:v1.0 .
docker push swr.cn-north-4.myhuaweicloud.com/your-namespace/zhitou-prediction:v1.0

# 3. 在CCE控制台部署
```

## 🔧 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `FLASK_ENV` | `development` | 运行环境 (development/production) |
| `PORT` | `9000` | 服务端口 |
| `HOST` | `0.0.0.0` | 监听地址 |
| `PYTHONUNBUFFERED` | `1` | Python输出缓冲 |

## 📊 监控和健康检查

### 健康检查端点

```bash
# 服务状态检查
curl http://localhost:9000/api/status

# 返回示例
{
  "status": "running",
  "ml_available": true,
  "supported_indices": 4,
  "timestamp": "2025-01-01T12:00:00.000000"
}
```

### Prometheus监控 (可选)

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'zhitou-prediction'
    static_configs:
      - targets: ['zhitou-prediction-service:80']
    metrics_path: '/api/status'
```

## 🚀 性能优化

### 资源配置建议

```yaml
# 生产环境资源配置
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

### 水平扩展

```bash
# 扩展副本数
kubectl scale deployment zhitou-prediction --replicas=5

# 自动扩展 (HPA)
kubectl autoscale deployment zhitou-prediction --cpu-percent=70 --min=2 --max=10
```

## 🔒 安全配置

### 网络策略

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: zhitou-prediction-netpol
spec:
  podSelector:
    matchLabels:
      app: zhitou-prediction
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: nginx-ingress
    ports:
    - protocol: TCP
      port: 9000
```

### HTTPS配置

```yaml
# ingress-tls.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: zhitou-prediction-ingress-tls
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - zhitou.yourdomain.com
    secretName: zhitou-tls
  rules:
  - host: zhitou.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: zhitou-prediction-service
            port:
              number: 80
```

## 🔧 故障排查

### 常见问题

1. **容器启动失败**
   ```bash
   docker logs zhitou-prediction
   kubectl describe pod <pod-name>
   ```

2. **服务无法访问**
   ```bash
   kubectl get svc
   kubectl describe svc zhitou-prediction-service
   ```

3. **健康检查失败**
   ```bash
   kubectl get pods
   kubectl logs <pod-name>
   ```

### 日志查看

```bash
# Docker
docker logs -f zhitou-prediction

# Kubernetes
kubectl logs -l app=zhitou-prediction -f --tail=100
```

## 📞 技术支持

如有部署问题，请检查：
1. 镜像是否正确构建
2. 网络端口是否开放
3. 环境变量是否正确设置
4. 资源配额是否充足

---

**注意**: 本部署指南适用于生产环境，请根据实际情况调整配置参数。
