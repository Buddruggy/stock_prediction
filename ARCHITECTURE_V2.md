# gogotou V2 - 架构设计（真实实现版）

> 本文档描述**当前真实代码**，与旧版（Go + 伪造历史 + LLM 编数）已完全不同。

## 设计原则

- **数据真实**：所有历史 K 线来自腾讯财经真实接口，绝无随机造假。
- **预测可回测**：价格预测由机器学习模型完成，置信度来自 walk-forward 回测的真实准确率，而非 LLM 自报。
- **LLM 只解读**：DeepSeek 仅把"真实指标 + 模型预测"转成中文解读，并做方向交叉验证；无 Key 时自动降级为纯模型模式。
- **零运维存储**：SQLite 文件即可运行，无需 MySQL。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Python 3.13 / FastAPI / uvicorn |
| 数据获取 | 腾讯财经 `web.ifzq.gtimg.cn` 日 K 线 |
| 特征/模型 | pandas / numpy / scikit-learn（GradientBoosting + RandomForest 集成） |
| 解读 | DeepSeek `deepseek-chat`（仅自然语言） |
| 存储 | SQLite |
| 调度 | APScheduler（工作日上海时间 15:10） |
| 前端 | Vue 3 + Chart.js + Nginx |

## 模块划分

```
backend/app/
├── main.py          FastAPI 入口 + 全部路由 + 启动调度器
├── config.py        环境变量配置（含 DEEPSEEK_API_KEY、端口、窗口大小）
├── data.py          腾讯真实 K 线获取（带内存缓存 + 容错重试）
├── indicators.py    真实技术指标：MA / RSI / MACD / 布林带 / 动量 / 波动率 / 量比
├── model.py         ML 预测 + walk-forward 回测（核心）
├── llm.py           DeepSeek 解读（仅自然语言，Key 走环境变量）
├── db.py            SQLite 持久化（预测记录 + 回测结果）
└── service.py       编排层：数据→指标→模型→LLM→存储
```

## 预测数据流

```
请求 / 定时(15:10)
   │
   ▼
data.py  ──► 腾讯真实日K线 (≈850 交易日, 带缓存)
   │
   ▼
indicators.py ──► 真实因子(MA/RSI/MACD/布林/动量/波动率/量比)
   │
   ▼
model.py
   ├─ 训练: 滚动窗口(默认250日) GradientBoosting + RandomForest
   ├─ 预测: 次日收益率 + 涨跌概率(prob_up)
   └─ 回测: walk-forward → 方向准确率 / 近期准确率 / MAE / RMSE
   │
   ▼
llm.py (可选)
   └─ 真实指标 + 模型预测 → DeepSeek 中文解读 + 独立方向校验
   │
   ▼
service.py ──► 统一方向信号(direction 与 direction_prob 永远一致)
   │
   ├─► db.py 写入 SQLite
   └─► 返回 JSON (current/predicted/changePercent/direction/confidence/backtest/...)
```

## 关键实现说明

### 1. 真实数据获取 (`data.py`)
- 接口：`https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?param={symbol},day,{start},{end},{limit},qfq`
- 单次最多返回约 850+ 条真实日 K；带内存缓存（`cache_ttl_seconds`）。
- 容错：请求失败指数退避重试，返回真实 OHLCV。

### 2. 模型与回测 (`model.py`)
- 特征：真实技术指标 + 收益率 + 价相对于均线的偏离。
- 训练：滚动窗口（默认 250 交易日）walk-forward，避免未来函数泄露。
- 预测头：①回归（次日收益率）→ 预测价；②分类（涨跌）→ prob_up。
- 回测：用历史窗口逐期预测，统计方向准确率、MAE、RMSE、样本数。
- 方向一致性：`direction` 由 `prob_up≥0.5` 推导，保证与 `direction_prob` 永不矛盾。

### 3. DeepSeek 解读 (`llm.py`)
- 仅当 `DEEPSEEK_API_KEY` 配置时启用，否则返回 `None`，系统降级为纯模型预测。
- 输入真实指标 + 模型预测 + 回测准确率，输出 `{reasoning, llm_direction, llm_confidence}`。
- 解析带 JSON 容错（兼容模型输出包裹文本/中文方向词）。

### 4. 定时任务 (`scheduler.py`)
- APScheduler：`CronTrigger(hour=15, minute=10, day_of_week='mon-fri', tz=Asia/Shanghai)`。
- 收盘后重算全部指数预测，并校验前一日预测方向是否正确（真实回测闭环）。

## API 契约（前端消费）

| 路径 | 返回关键字段 |
| --- | --- |
| `/api/v1/indices` | 指数元信息 |
| `/api/v1/predict/{code}` | `current/predicted/change/changePercent/predicted_return/prob_up/direction/confidence/technical_indicators/backtest/feature_importance/reasoning/llm_direction` |
| `/api/v1/predict/all` | 各指数预测 dict |
| `/api/v1/history/{code}?period=` | 真实 K 线数组（date/open/close/high/low/volume） |
| `/api/v1/predict/history/all` | 历史预测记录（按 code 分组） |
| `/api/v1/prediction-stats` | `total_predictions/success_rate/avg_confidence` |
| `/api/v1/health` | 健康检查 |

## 部署

- `docker-compose.yml`：backend(Python, 8001) + frontend(Vue/Nginx, 80/9000/443)，SQLite 卷持久化，无 MySQL 依赖。
- 旧 `backend-go/`（Go 实现）已废弃，不再被任何构建引用。
