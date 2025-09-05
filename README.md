# 智投预测 - AI股市指数预测平台

<div align="center">
  <img src="frontend/public/logo.svg" alt="智投预测Logo" width="200">
  
  **基于人工智能的中国股票指数预测平台**
  
  [![Docker](https://img.shields.io/badge/Docker-支持-blue?logo=docker)](https://www.docker.com/)
  [![Python](https://img.shields.io/badge/Python-3.8+-green?logo=python)](https://www.python.org/)
  [![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
  [![AI](https://img.shields.io/badge/AI-机器学习-red?logo=tensorflow)](https://tensorflow.org/)
</div>

---

一个现代化的股票指数预测网站，专注于预测中国主要股票指数（上证综指、深证成指、创业板指、科创50）的明日涨跌情况。

## 🌟 功能特性

### 📊 指数预测
- **上证综指 (SH000001)** - 上海证券交易所综合股价指数
- **深证成指 (SZ399001)** - 深圳证券交易所成分股指数
- **创业板指 (SZ399006)** - 创业板综合指数
- **科创50 (SH000688)** - 科创板50成分指数

### 🤖 智能预测模型
- **机器学习预测**: 使用随机森林回归模型，结合技术指标进行预测
- **技术指标分析**: 包含移动平均线、RSI、MACD、布林带等多种技术指标
- **置信度评估**: 为每个预测提供置信度评分
- **实时数据**: 自动获取最新的股票数据进行分析

### 💻 现代化界面
- **响应式设计**: 完美适配桌面和移动设备
- **实时更新**: 自动更新市场数据和预测结果
- **交互图表**: 使用Chart.js展示历史趋势和预测数据
- **直观展示**: 清晰的卡片式布局展示预测信息

## 🚀 快速开始

### 方式一: Docker部署 (推荐)

**环境要求**: Docker 20.0+ 和 Docker Compose 2.0+

```bash
# 1. 克隆项目
git clone <项目地址>
cd stock_prediction

# 2. 一键启动 (推荐)
./scripts/docker-run.sh

# 3. 访问网站
# http://localhost:9000
```

### 方式二: 本地开发环境

**环境要求**: Python 3.8+ 和现代浏览器

```bash
# 1. 克隆项目
git clone <项目地址>
cd stock_prediction

# 2. 创建虚拟环境
python -m venv venv
source venv/bin/activate  # Linux/Mac
# 或 venv\\Scripts\\activate  # Windows

# 3. 安装依赖
pip install -r requirements.txt

# 4. 启动服务器
python app.py

# 5. 访问网站
# http://localhost:9000
```

## 📁 项目结构

```
stock_prediction/
├── app.py              # Flask后端服务器
├── index.html          # 前端主页面
├── styles.css          # 样式文件
├── script.js           # JavaScript交互逻辑
├── requirements.txt    # Python依赖
└── README.md          # 项目说明
```

## 🔧 API 接口

### 获取支持的指数列表
```http
GET /api/indices
```

### 预测指定指数
```http
GET /api/predict/<index_code>
```
参数: `index_code` - 指数代码 (sh000001, sz399001, sz399006, sh000688)

### 预测所有指数
```http
GET /api/predict/all
```

### 获取历史数据
```http
GET /api/history/<index_code>?period=1mo
```
参数: 
- `index_code` - 指数代码
- `period` - 时间周期 (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)

### 获取服务状态
```http
GET /api/status
```

## 📊 预测模型说明

### 数据来源
- 使用 `yfinance` 库获取实时股票数据
- 数据包含开盘价、最高价、最低价、收盘价、成交量等信息

### 技术指标
1. **移动平均线**: MA5、MA10、MA20
2. **相对强弱指数**: RSI (14日)
3. **MACD**: 指数平滑移动平均线
4. **布林带**: 20日布林带上下轨
5. **成交量指标**: 成交量比率

### 预测算法
- **主要模型**: 随机森林回归 (Random Forest Regressor)
- **备用模型**: 线性趋势预测 (当机器学习库不可用时)
- **特征工程**: 使用过去30天的技术指标作为特征
- **置信度**: 基于历史预测准确率和价格波动性计算

## 🎨 界面特性

### 设计理念
- **现代扁平化**: 采用现代扁平化设计风格
- **渐变配色**: 使用优雅的渐变色彩方案
- **卡片布局**: 清晰的卡片式信息展示
- **动画效果**: 流畅的过渡和悬停动画

### 响应式适配
- **桌面端**: 多栏网格布局，充分利用屏幕空间
- **平板端**: 自适应两栏布局
- **手机端**: 单栏垂直布局，优化触摸操作

### 交互体验
- **平滑滚动**: 导航链接支持平滑滚动
- **模态框**: 详细信息以模态框形式展示
- **实时更新**: 数据自动刷新，无需手动操作
- **键盘支持**: ESC键关闭模态框等快捷操作

## ⚙️ 配置选项

### 数据缓存
- 默认缓存时间: 5分钟
- 可在 `app.py` 中修改 `cache_duration` 参数

### 预测参数
- 历史数据窗口: 30天 (可在 `prepare_features` 函数中调整)
- 模型训练比例: 80% (可在 `ml_predict` 函数中调整)

## 🚨 注意事项

### 免责声明
- **仅供参考**: 本系统提供的预测结果仅供参考，不构成投资建议
- **投资风险**: 股票投资存在风险，投资需谨慎
- **数据延迟**: 数据可能存在延迟，请以实际市场数据为准

### 使用限制
- **数据来源**: 依赖第三方数据源，可能受到网络和服务限制
- **预测准确性**: 预测结果基于历史数据，无法保证未来准确性
- **服务可用性**: 服务器可能因维护或其他原因暂时不可用

## 🔄 更新日志

### v1.0.0 (2024-01-01)
- ✨ 初始版本发布
- 🎯 支持4个主要中国股票指数预测
- 🤖 集成机器学习预测模型
- 💻 现代化响应式界面
- 📊 实时数据获取和展示
- 📱 移动端适配

## 🤝 贡献指南

欢迎提交问题和改进建议！

### 开发环境设置
1. Fork 项目
2. 创建功能分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送分支: `git push origin feature/new-feature`
5. 提交 Pull Request

### 代码规范
- Python代码遵循 PEP 8 规范
- JavaScript代码使用 ES6+ 语法
- CSS使用BEM命名规范
- 提交信息使用约定式提交格式

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📞 联系方式

如有问题或建议，请通过以下方式联系:
- 📧 邮箱: [your-email@example.com]
- 🐛 问题反馈: [GitHub Issues]
- 💬 讨论: [GitHub Discussions]

---

**⚠️ 风险提示**: 本系统仅供学习和研究使用，预测结果不构成投资建议。股市有风险，投资需谨慎！
