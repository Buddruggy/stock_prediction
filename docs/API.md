# 智投预测 - API文档

## 📚 API接口说明

### 基础信息

- **基础URL**: `http://localhost:9000`
- **响应格式**: JSON
- **编码**: UTF-8

### 🔍 状态接口

#### GET /api/status

获取服务状态信息

**响应示例:**
```json
{
  "status": "running",
  "ml_available": true,
  "supported_indices": 4,
  "timestamp": "2025-01-01T12:00:00.000000",
  "cache_duration": 300
}
```

### 📊 指数接口

#### GET /api/indices

获取支持的股票指数列表

**响应示例:**
```json
{
  "sh000001": {
    "name": "上证综合指数",
    "symbol": "000001.SS",
    "code": "SH000001"
  },
  "sz399001": {
    "name": "深证成分指数", 
    "symbol": "399001.SZ",
    "code": "SZ399001"
  }
}
```

### 🤖 预测接口

#### GET /api/predict/{index_code}

获取指定指数的预测数据

**参数:**
- `index_code`: 指数代码 (sh000001, sz399001, sz399006, sh000688)

**响应示例:**
```json
{
  "code": "sh000001",
  "name": "上证综合指数",
  "symbol": "SH000001",
  "current": 3245.67,
  "change": 12.34,
  "changePercent": 0.38,
  "predicted": 3268.45,
  "predictedChange": 22.78,
  "predictedPercent": 0.70,
  "confidence": 78.0,
  "timestamp": "2025-01-01T12:00:00.000000"
}
```

#### GET /api/predict/all

获取所有指数的预测数据

**响应示例:**
```json
{
  "sh000001": {
    "code": "sh000001",
    "name": "上证综合指数",
    "current": 3245.67,
    "predicted": 3268.45,
    "confidence": 78.0
  },
  "sz399001": {
    ...
  }
}
```

### 📈 历史数据接口

#### GET /api/history/{index_code}

获取指定指数的历史数据

**参数:**
- `index_code`: 指数代码
- `period`: 时间周期 (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)

**示例:** `/api/history/sh000001?period=1mo`

**响应示例:**
```json
{
  "code": "sh000001",
  "name": "上证综合指数",
  "history": [
    {
      "date": "2025-01-01",
      "open": 3240.12,
      "high": 3250.45,
      "low": 3235.67,
      "close": 3245.67,
      "volume": 123456789
    }
  ]
}
```

### ❌ 错误响应

所有错误响应都包含错误信息：

```json
{
  "error": "错误描述信息"
}
```

**常见状态码:**
- `200`: 成功
- `400`: 请求参数错误
- `404`: 资源不存在
- `500`: 服务器内部错误

### 🔒 使用限制

- 请求频率限制: 100次/小时
- 单次请求超时: 30秒
- 数据更新频率: 5分钟

### 💡 使用建议

1. 建议在生产环境中实现请求缓存
2. 处理网络超时和错误情况
3. 预测结果仅供参考，不构成投资建议
