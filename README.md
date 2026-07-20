# gogotou - AI 股市指数预测平台

<div align="center">
  <img src="frontend/public/logo.svg" alt="gogotou Logo" width="200">

  **基于真实行情数据 + 机器学习模型的中国股票指数预测平台**

  [![Python](https://img.shields.io/badge/Python-3.13-blue?logo=python)](https://python.org/)
  [![FastAPI](https://img.shields.io/badge/FastAPI-后端-green?logo=fastapi)](https://fastapi.tiangolo.com/)
  [![scikit-learn](https://img.shields.io/badge/ml-scikit--learn-orange?logo=scikit-learn)](https://scikit-learn.org/)
  [![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
</div>

---

预测 **中国主要股票指数**（上证综指、深证成指、创业板指、科创 50）的**次日涨跌方向与幅度**。

## ⚠️ 本版本与旧版的本质区别

旧版（Go）的"预测"名不副实：**历史数据是随机数伪造的，预测靠通用大模型在 ±5% 内编一个数，且硬编码了一个已泄露的 API Key**。本版本是完全重写：

| 维度 | 旧版（Go） | 本版本（Python） |
| --- | --- | --- |
| 历史数据 | `rand` 随机生成（造假） | **腾讯财经真实日 K 线**（约 850 交易日） |
| 预测引擎 | 通用 LLM 编数 | **集成 ML 模型**（GradientBoosting + 随机森林投票） |
| 置信度 | LLM 自报（不可信） | **walk-forward 回测的真实方向准确率** |
| DeepSeek | 既预测又解读，Key 硬编码 | **仅做自然语言解读**，Key 走环境变量 |
| 存储 | MySQL（重） | **SQLite**（零配置，可直接跑） |

> 诚实提示：纯技术指标预测指数日线，长期方向准确率在真实市场里通常接近随机（本系统回测约 50%–60%）。这不是"稳赚信号"，而是**可回测、可解释、有诚实置信度**的研究工具。

## 🌟 功能特性

- **4 大指数**：上证综指 (sh000001)、深证成指 (sz399001)、创业板指 (sz399006)、科创 50 (sh000688)
- **真实数据**：腾讯财经 `web.ifzq.gtimg.cn` 真实日 K 线，含 OHLCV
- **真模型预测**：基于真实指标训练 GradientBoosting/RandomForest，预测次日收益率与涨跌方向
- **真实回测**：walk-forward 滚动窗口回测，给出方向准确率、近期准确率、MAE、RMSE
- **DeepSeek 解读**：输入真实指标 + 模型预测，生成中文解读与独立方向交叉验证（可选）
- **定时任务**：每个工作日上海时间 15:10（收盘后）自动重算，并校验前一日预测方向
- **现代前端**：Vue 3 + Chart.js，响应式卡片展示

## 🚀 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 1. 配置 DeepSeek Key（可选，留空则纯模型模式）
cp .env.example .env
# 编辑 .env，填入 DEEPSEEK_API_KEY

# 2. 一键启动
docker compose up -d --build

# 访问
#   前端: http://localhost:80
#   后端 API: http://localhost:8001
#   API 文档: http://localhost:8001/docs
```

### 方式二：本地开发

**环境要求**：Python 3.13+、Node 18+

```bash
# 后端（终端 1）
python3 -m venv backend/.venv && source backend/.venv/bin/activate
pip install -r backend/requirements.txt
cd backend && python -m uvicorn app.main:app --host 0.0.0.0 --port 8001 --reload

# 前端（终端 2）
cd frontend && npm install && npm run dev
# 前端默认代理到 http://localhost:8001
```

## 📁 项目结构

```
stock_prediction/
├── backend/                 # Python / FastAPI 后端
│   ├── app/
│   │   ├── main.py          # FastAPI 入口 + 路由
│   │   ├── config.py        # 环境变量配置
│   │   ├── data.py          # 腾讯真实 K 线获取 + 缓存
│   │   ├── indicators.py    # 真实技术指标计算
│   │   ├── model.py         # ML 预测 + walk-forward 回测
│   │   ├── llm.py           # DeepSeek 解读（仅自然语言）
│   │   ├── db.py            # SQLite 存储
│   │   └── service.py       # 核心编排（数据→指标→模型→LLM）
│   ├── requirements.txt
│   └── predictions.db       # SQLite 数据（运行时生成）
├── frontend/                # Vue 3 前端
├── deployment/              # Dockerfile（backend / frontend）
├── docker-compose.yml
└── Makefile
```

## 🔧 API 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v1/indices` | 支持的指数列表 |
| GET | `/api/v1/predict/{code}` | 预测指定指数（实时计算，首次会训练模型） |
| GET | `/api/v1/predict/all` | 预测全部指数 |
| GET | `/api/v1/history/{code}?period=1mo` | K 线历史（腾讯真实数据） |
| GET | `/api/v1/predict/history/all` | 历史预测记录 |
| GET | `/api/v1/prediction-stats` | 预测统计（成功率/平均置信度） |
| GET | `/api/v1/health` | 健康检查 |

## 📊 预测原理（诚实版）

1. **真实数据**：从腾讯财经拉取约 850 个交易日的真实日 K 线。
2. **真实指标**：基于真实数据计算 MA5/10/20、RSI(14)、MACD、布林带、动量、波动率、量比等因子。
3. **ML 预测**：用滚动窗口（默认 250 交易日）训练集成模型，预测**次日收益率**与**涨跌概率**。
4. **回测校验**：walk-forward 回测给出真实的历史方向准确率（替代旧版"LLM 自报 90% 置信度"）。
5. **DeepSeek 解读**（可选）：把真实指标 + 模型预测喂给 DeepSeek，生成中文解读，并对其方向做独立交叉验证。

## ⚙️ 配置项（环境变量）

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `DEEPSEEK_API_KEY` | DeepSeek Key（仅解读用，留空=纯模型） | 空 |
| `PORT` | 后端端口 | 8001 |
| `DB_PATH` | SQLite 路径 | `./backend/predictions.db` |
| `HISTORY_DAYS` | 拉取的真实交易日数 | 400 |
| `TRAIN_WINDOW` | 滚动训练窗口 | 250 |
| `PREDICT_HOUR` / `PREDICT_MINUTE` | 定时预测时间（上海时区） | 15 / 10 |

## 🚨 免责声明

- **仅供参考**：预测结果基于历史数据统计，不构成任何投资建议。
- **投资风险**：股市有风险，决策需谨慎，亏损自负。
- **数据延迟**：依赖第三方数据源，可能存在延迟或中断。

## 📄 许可证

MIT — 详见 [LICENSE](LICENSE)。
