<template>
  <div class="modern-history">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">历史预测记录</h1>
      <p class="page-subtitle">查看和分析历史股票预测数据</p>
    </div>

    <!-- 控制面板 -->
    <div class="control-panel">
      <div class="filter-section">
        <div class="filter-item">
          <label class="filter-label">时间范围:</label>
          <select v-model="selectedDays" @change="fetchHistoricalData" class="filter-select">
            <option value="7">近7天</option>
            <option value="15">近15天</option>
            <option value="30">近30天</option>
            <option value="60">近60天</option>
            <option value="90">近90天</option>
          </select>
        </div>
        
        <button @click="fetchHistoricalData" :disabled="loading" class="refresh-button">
          <span v-if="!loading">🔄</span>
          <span v-else>⏳</span>
          刷新数据
        </button>
      </div>
    </div>

    <!-- 统计概览 -->
    <div v-if="hasData" class="stats-overview">
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-number">{{ totalPredictions }}</div>
          <div class="stat-label">预测记录总数</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">🎯</div>
        <div class="stat-content">
          <div class="stat-number">{{ avgConfidence }}%</div>
          <div class="stat-label">平均置信度</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📈</div>
        <div class="stat-content">
          <div class="stat-number">{{ Object.keys(historicalData).length }}</div>
          <div class="stat-label">覆盖指数数量</div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="status-section loading">
      <div class="status-card">
        <div class="loading-spinner"></div>
        <p class="status-text">正在加载历史预测数据...</p>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="status-section error">
      <div class="status-card">
        <div class="error-icon">⚠️</div>
        <p class="status-text">{{ error }}</p>
        <button @click="fetchHistoricalData" class="retry-button">重试</button>
      </div>
    </div>

    <!-- 空数据状态 -->
    <div v-else-if="!hasData" class="status-section empty">
      <div class="status-card">
        <div class="empty-icon">📊</div>
        <p class="status-text">暂无历史预测数据</p>
        <p class="empty-description">请调整筛选条件或稍后再试</p>
      </div>
    </div>

    <!-- 历史数据图表 -->
    <div v-else class="history-charts">
      <div v-for="(predictions, indexCode) in historicalData" :key="indexCode" class="chart-section">
        <div class="chart-header">
          <h3 class="chart-title">{{ getIndexName(indexCode) }}</h3>
          <div class="chart-stats">
            <span class="stat-item">记录数: {{ predictions.length }}</span>
            <span class="stat-item">平均置信度: {{ getAvgConfidence(predictions) }}%</span>
          </div>
        </div>
        
        <div class="chart-container">
          <canvas 
            :ref="el => setChartRef(el, indexCode)" 
            class="price-chart"
          ></canvas>
          
          <!-- 备用显示：如果图表加载失败 -->
          <div v-if="!charts[indexCode]" class="chart-fallback">
            <div class="fallback-message">
              <p>图表加载中...</p>
            </div>
          </div>
        </div>
        
        <!-- 图表说明 -->
        <div class="chart-legend">
          <div class="legend-item">
            <div class="legend-color current"></div>
            <span>当前价格</span>
          </div>
          <div class="legend-item">
            <div class="legend-color predicted"></div>
            <span>预测价格</span>
          </div>
        </div>
        
        <!-- 备用数据表格（如果图表失败） -->
        <div v-if="showFallbackTable" class="fallback-table">
          <h4>数据概览</h4>
          <div class="simple-data-list">
            <div v-for="(prediction, index) in predictions.slice(0, 5)" :key="index" class="data-item">
              <span class="date">{{ formatDate(prediction.timestamp || prediction.prediction_date) }}</span>
              <span class="current-price">{{ prediction.current?.toFixed(2) || '--' }}</span>
              <span class="predicted-price">{{ prediction.predicted?.toFixed(2) || '--' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 技术指标模态框 -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h4 class="modal-title">技术指标详情</h4>
          <button @click="closeModal" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <div v-if="selectedPrediction?.technical_indicators" class="indicator-grid">
            <div v-for="(value, key) in selectedPrediction.technical_indicators" :key="key" class="indicator-item">
              <div class="indicator-label">{{ key }}</div>
              <div class="indicator-value">{{ typeof value === 'number' ? value.toFixed(4) : value }}</div>
            </div>
            
            <div class="indicator-item">
              <div class="indicator-label">预测日期</div>
              <div class="indicator-value">{{ formatDate(selectedPrediction.timestamp) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import axios from 'axios'

// 响应式数据
const loading = ref(false)
const error = ref('')
const historicalData = ref({})
const selectedIndex = ref('all')
const selectedDays = ref(30)
const showModal = ref(false)
const selectedPrediction = ref(null)
const charts = ref({})
const showFallbackTable = ref(false)

// 设置图表引用
const setChartRef = (el, indexCode) => {
  if (el) {
    charts.value[indexCode] = el
  }
}

// 指数名称映射
const indexNames = {
  'sh000001': '上证综指',
  'sz399001': '深证成指',
  'sz399006': '创业板指',
  'sh000688': '科创50'
}

// 计算属性
const hasData = computed(() => {
  return Object.keys(historicalData.value).length > 0
})

const totalPredictions = computed(() => {
  return Object.values(historicalData.value).reduce((total, predictions) => {
    return total + (predictions?.length || 0)
  }, 0)
})

const avgConfidence = computed(() => {
  let totalConfidence = 0
  let count = 0
  
  Object.values(historicalData.value).forEach(predictions => {
    predictions?.forEach(prediction => {
      if (prediction.confidence) {
        totalConfidence += prediction.confidence
        count++
      }
    })
  })
  
  return count > 0 ? (totalConfidence / count).toFixed(1) : '0'
})

// 方法
const fetchHistoricalData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    // 始终获取所有指数的数据
    const url = `/api/v1/predict/history/all?days=${selectedDays.value}`
    
    const response = await axios.get(url)
    
    if (response.data.code === 200) {
      historicalData.value = response.data.data || {}
    } else {
      error.value = response.data.message || '获取数据失败'
      historicalData.value = {}
    }
  } catch (err) {
    console.error('获取历史预测数据失败:', err)
    error.value = '网络错误，请检查网络连接后重试'
    historicalData.value = {}
  } finally {
    loading.value = false
  }
}

const getIndexName = (indexCode) => {
  return indexNames[indexCode] || indexCode
}

const formatDate = (timestamp) => {
  if (!timestamp) return '--'
  // 如果是日期格式 (YYYY-MM-DD)直接返回
  if (timestamp.length === 10 && timestamp.includes('-')) {
    return timestamp
  }
  // 否则解析为日期
  const date = new Date(timestamp)
  return date.toISOString().split('T')[0]
}

const formatChange = (change) => {
  if (!change && change !== 0) return '--'
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(2)}`
}

const formatPercent = (percent) => {
  if (!percent && percent !== 0) return '--'
  const sign = percent >= 0 ? '+' : ''
  return `${sign}${percent.toFixed(2)}%`
}

const getChangeClass = (change) => {
  if (!change && change !== 0) return ''
  return change >= 0 ? 'positive' : 'negative'
}

const getConfidenceClass = (confidence) => {
  if (!confidence) return 'low'
  if (confidence >= 80) return 'high'
  if (confidence >= 60) return 'medium'
  return 'low'
}

const getAvgConfidence = (predictions) => {
  if (!predictions || predictions.length === 0) return '0'
  
  const total = predictions.reduce((sum, pred) => {
    return sum + (pred.confidence || 0)
  }, 0)
  
  return (total / predictions.length).toFixed(1)
}

const showIndicators = (prediction) => {
  selectedPrediction.value = prediction
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedPrediction.value = null
}

// 处理图表数据：预测价格后移一天
const processChartData = (predictions) => {
  if (!predictions || predictions.length === 0) return { labels: [], currentPrices: [], predictedPrices: [] }
  
  // 按日期排序
  const sortedPredictions = [...predictions].sort((a, b) => {
    const dateA = new Date(a.timestamp || a.prediction_date)
    const dateB = new Date(b.timestamp || b.prediction_date)
    return dateA - dateB
  })
  
  const labels = []
  const currentPrices = []
  const predictedPrices = []
  
  sortedPredictions.forEach((prediction, index) => {
    const date = formatDate(prediction.timestamp || prediction.prediction_date)
    labels.push(date)
    currentPrices.push(prediction.current || 0)
    
    // 预测价格前移一天：今天的预测价格展示在明天的时间轴上
    if (index > 0) {
      const prevPrediction = sortedPredictions[index - 1]
      predictedPrices.push(prevPrediction.predicted || 0)
    } else {
      // 第一个预测没有对应的前一天预测数据
      predictedPrices.push(null)
    }
  })
  
  return { labels, currentPrices, predictedPrices }
}

// 绘制折线图
const drawChart = (canvas, predictions) => {
  if (!canvas || !predictions || predictions.length === 0) return
  
  try {
    const ctx = canvas.getContext('2d')
    const { labels, currentPrices, predictedPrices } = processChartData(predictions)
    
    // 清除画布
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    
    if (labels.length === 0) return
  
    // 设置画布尺寸
    const rect = canvas.getBoundingClientRect()
    canvas.width = rect.width * window.devicePixelRatio
    canvas.height = rect.height * window.devicePixelRatio
    ctx.scale(window.devicePixelRatio, window.devicePixelRatio)
    
    // 检测移动端
    const isMobile = window.innerWidth <= 768
    const isSmallMobile = window.innerWidth <= 480
    
    // 根据屏幕尺寸调整参数
    const padding = isMobile ? (isSmallMobile ? 20 : 30) : 40
    const chartWidth = rect.width - padding * 2
    const chartHeight = rect.height - padding * 2
  
  // 计算价格范围
  const allPrices = [...currentPrices, ...predictedPrices.filter(p => p !== null && p > 0)]
  if (allPrices.length === 0) return
  
  const minPrice = Math.min(...allPrices)
  const maxPrice = Math.max(...allPrices)
  const priceRange = maxPrice - minPrice
  const pricePadding = priceRange > 0 ? priceRange * 0.1 : maxPrice * 0.1
  
  // 绘制网格
  ctx.strokeStyle = '#e5e7eb'
  ctx.lineWidth = 1
  
  // 水平网格线 - 移动端减少数量
  const gridYCount = isMobile ? (isSmallMobile ? 3 : 4) : 5
  for (let i = 0; i <= gridYCount; i++) {
    const y = padding + (chartHeight / gridYCount) * i
    ctx.beginPath()
    ctx.moveTo(padding, y)
    ctx.lineTo(padding + chartWidth, y)
    ctx.stroke()
  }
  
  // 垂直网格线 - 移动端减少数量
  const gridXCount = isMobile ? Math.max(2, Math.floor(labels.length / 2)) : labels.length - 1
  for (let i = 0; i <= gridXCount; i++) {
    const x = padding + (chartWidth / gridXCount) * i
    ctx.beginPath()
    ctx.moveTo(x, padding)
    ctx.lineTo(x, padding + chartHeight)
    ctx.stroke()
  }
  
  // 绘制当前价格线
  ctx.strokeStyle = '#3b82f6'
  ctx.lineWidth = 3
  ctx.beginPath()
  
  currentPrices.forEach((price, index) => {
    const x = padding + (chartWidth / (labels.length - 1)) * index
    const y = padding + chartHeight - ((price - minPrice + pricePadding) / (priceRange + pricePadding * 2)) * chartHeight
    
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()
  
  // 绘制预测价格线
  ctx.strokeStyle = '#ef4444'
  ctx.lineWidth = 3
  ctx.setLineDash([5, 5])
  ctx.beginPath()
  
  predictedPrices.forEach((price, index) => {
    if (price !== null) {
      const x = padding + (chartWidth / (labels.length - 1)) * index
      const y = padding + chartHeight - ((price - minPrice + pricePadding) / (priceRange + pricePadding * 2)) * chartHeight
      
      if (index === 0) {
        ctx.moveTo(x, y)
      } else {
        ctx.lineTo(x, y)
      }
    }
  })
  ctx.stroke()
  ctx.setLineDash([])
  
  // 绘制数据点
  ctx.fillStyle = '#3b82f6'
  currentPrices.forEach((price, index) => {
    const x = padding + (chartWidth / (labels.length - 1)) * index
    const y = padding + chartHeight - ((price - minPrice + pricePadding) / (priceRange + pricePadding * 2)) * chartHeight
    
    ctx.beginPath()
    ctx.arc(x, y, 4, 0, 2 * Math.PI)
    ctx.fill()
  })
  
  ctx.fillStyle = '#ef4444'
  predictedPrices.forEach((price, index) => {
    if (price !== null) {
      const x = padding + (chartWidth / (labels.length - 1)) * index
      const y = padding + chartHeight - ((price - minPrice + pricePadding) / (priceRange + pricePadding * 2)) * chartHeight
      
      ctx.beginPath()
      ctx.arc(x, y, 4, 0, 2 * Math.PI)
      ctx.fill()
    }
  })
  
  // 绘制Y轴标签
  ctx.fillStyle = '#6b7280'
  ctx.font = isMobile ? (isSmallMobile ? '10px sans-serif' : '11px sans-serif') : '12px sans-serif'
  ctx.textAlign = 'right'
  
  // 移动端减少Y轴标签数量
  const yLabelCount = isMobile ? (isSmallMobile ? 3 : 4) : 5
  
  for (let i = 0; i <= yLabelCount; i++) {
    const price = maxPrice - (priceRange / yLabelCount) * i
    const y = padding + (chartHeight / yLabelCount) * i + 4
    
    // 移动端简化价格显示
    let priceText
    if (isSmallMobile) {
      priceText = price.toFixed(0) // 小屏只显示整数
    } else if (isMobile) {
      priceText = price.toFixed(1) // 中屏显示一位小数
    } else {
      priceText = price.toFixed(2) // 大屏显示两位小数
    }
    
    ctx.fillText(priceText, padding - 5, y)
  }
  
  // 绘制X轴标签
  ctx.textAlign = 'center'
  ctx.font = isMobile ? (isSmallMobile ? '9px sans-serif' : '10px sans-serif') : '12px sans-serif'
  
  // 移动端减少X轴标签数量，避免重叠
  const xLabelStep = isMobile ? Math.max(1, Math.floor(labels.length / (isSmallMobile ? 3 : 4))) : 1
  
  labels.forEach((label, index) => {
    // 只在指定间隔显示标签
    if (index % xLabelStep === 0 || index === labels.length - 1) {
      const x = padding + (chartWidth / (labels.length - 1)) * index
      const y = padding + chartHeight + (isMobile ? 15 : 20)
      
      // 移动端简化日期显示
      let labelText
      if (isSmallMobile) {
        // 小屏只显示月-日
        labelText = label.split('-').slice(1).join('-')
      } else if (isMobile) {
        // 中屏显示月-日
        labelText = label.split('-').slice(1).join('-')
      } else {
        // 大屏显示完整日期
        labelText = label
      }
      
      ctx.fillText(labelText, x, y)
    }
  })
  
  } catch (error) {
    console.error('绘制图表时出错:', error)
  }
}

// 初始化所有图表
const initCharts = async () => {
  try {
    await nextTick()
    
    console.log('初始化图表，数据:', historicalData.value)
    
    let hasValidCharts = false
    
    Object.keys(historicalData.value).forEach(indexCode => {
      const canvas = charts.value[indexCode]
      const predictions = historicalData.value[indexCode]
      
      console.log(`处理指数 ${indexCode}:`, { canvas: !!canvas, predictionsCount: predictions?.length })
      
      if (canvas && predictions && predictions.length > 0) {
        try {
          drawChart(canvas, predictions)
          hasValidCharts = true
        } catch (error) {
          console.error(`绘制指数 ${indexCode} 图表失败:`, error)
        }
      } else {
        console.warn(`跳过指数 ${indexCode}:`, { canvas: !!canvas, predictionsCount: predictions?.length })
      }
    })
    
    // 如果所有图表都失败了，显示备用表格
    showFallbackTable.value = !hasValidCharts
    
  } catch (error) {
    console.error('初始化图表时出错:', error)
    showFallbackTable.value = true
  }
}

// 监听数据变化
watch(historicalData, () => {
  initCharts()
}, { deep: true })

// 组件挂载时获取数据
onMounted(() => {
  fetchHistoricalData()
})
</script>

<style lang="scss" scoped>
// 使用内联样式变量定义，避免外部依赖
:root {
  --claude-space: 8px;
  --claude-space-sm: 4px;
  --claude-space-lg: 16px;
  --claude-space-xl: 24px;
  --claude-space-xs: 2px;
  --claude-radius: 8px;
  --claude-radius-lg: 12px;
  --claude-border: #e5e7eb;
  --claude-bg-primary: #ffffff;
  --claude-bg-secondary: #f9fafb;
  --claude-bg-tertiary: #f3f4f6;
  --claude-text-primary: #111827;
  --claude-text-secondary: #6b7280;
  --claude-text-tertiary: #9ca3af;
  --claude-primary: #3b82f6;
  --claude-success: #10b981;
  --claude-warning: #f59e0b;
  --claude-danger: #ef4444;
  --claude-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  --claude-shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.modern-history {
  min-height: 100vh;
  padding: var(--claude-space-xl);
  background: var(--claude-bg-primary);
  
  @media (max-width: 768px) {
    padding: var(--claude-space-lg);
  }
}

.page-header {
  text-align: center;
  margin-bottom: var(--claude-space-xl);
  
  .page-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--claude-text-primary);
    margin-bottom: var(--claude-space);
    
    @media (max-width: 768px) {
      font-size: 2rem;
    }
  }
  
  .page-subtitle {
    font-size: 1.1rem;
    color: var(--claude-text-secondary);
    margin: 0;
  }
}

.control-panel {
  background: var(--claude-bg-primary);
  border: 1px solid var(--claude-border);
  border-radius: var(--claude-radius-lg);
  padding: var(--claude-space-xl);
  box-shadow: var(--claude-shadow);
  margin-bottom: var(--claude-space-xl);
  
  .filter-section {
    display: flex;
    align-items: center;
    gap: var(--claude-space-lg);
    flex-wrap: wrap;
    
    @media (max-width: 768px) {
      flex-direction: column;
      align-items: stretch;
    }
  }
  
  .filter-item {
    display: flex;
    align-items: center;
    gap: var(--claude-space);
    
    @media (max-width: 768px) {
      flex-direction: column;
      align-items: stretch;
    }
  }
  
  .filter-label {
    font-weight: 600;
    color: var(--claude-text-primary);
    white-space: nowrap;
  }
  
  .filter-select {
    padding: var(--claude-space) var(--claude-space-lg);
    border: 1px solid var(--claude-border);
    border-radius: var(--claude-radius);
    background: var(--claude-bg-secondary);
    color: var(--claude-text-primary);
    font-size: 0.9rem;
    min-width: 120px;
    
    &:focus {
      outline: none;
      border-color: var(--claude-primary);
    }
  }
  
  .refresh-button {
    background: var(--claude-primary);
    color: white;
    border: none;
    padding: var(--claude-space) var(--claude-space-lg);
    border-radius: var(--claude-radius);
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: var(--claude-space-sm);
    
    &:hover {
      background: #2563eb;
    }
    
    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.stats-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--claude-space-lg);
  margin-bottom: var(--claude-space-xl);
  
  .stat-card {
    background: var(--claude-bg-primary);
    border: 1px solid var(--claude-border);
    border-radius: var(--claude-radius-lg);
    padding: var(--claude-space-xl);
    box-shadow: var(--claude-shadow);
    display: flex;
    align-items: center;
    gap: var(--claude-space-lg);
    
    .stat-icon {
      font-size: 2rem;
      opacity: 0.8;
    }
    
    .stat-content {
      .stat-number {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--claude-primary);
        margin-bottom: var(--claude-space-xs);
      }
      
      .stat-label {
        font-size: 0.9rem;
        color: var(--claude-text-secondary);
      }
    }
  }
}

.history-charts {
  .chart-section {
    background: var(--claude-bg-primary);
    border: 1px solid var(--claude-border);
    border-radius: var(--claude-radius-lg);
    padding: 0;
    box-shadow: var(--claude-shadow);
    margin-bottom: var(--claude-space-xl);
    overflow: hidden;
    
    .chart-header {
      padding: var(--claude-space-lg);
      border-bottom: 1px solid var(--claude-border);
      background: var(--claude-bg-tertiary);
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: var(--claude-space);
      
      .chart-title {
        font-size: 1.2rem;
        font-weight: 600;
        color: var(--claude-text-primary);
        margin: 0;
      }
      
      .chart-stats {
        display: flex;
        gap: var(--claude-space-lg);
        
        .stat-item {
          font-size: 0.9rem;
          color: var(--claude-text-secondary);
        }
      }
    }
    
    .chart-container {
      padding: var(--claude-space-lg);
      position: relative;
      
      .price-chart {
        width: 100%;
        height: 400px;
        border-radius: var(--claude-radius);
        
        @media (max-width: 768px) {
          height: 300px;
        }
        
        @media (max-width: 480px) {
          height: 250px;
        }
      }
      
      .chart-fallback {
        position: absolute;
        top: var(--claude-space-lg);
        left: var(--claude-space-lg);
        right: var(--claude-space-lg);
        bottom: var(--claude-space-lg);
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--claude-bg-secondary);
        border-radius: var(--claude-radius);
        
        .fallback-message {
          text-align: center;
          color: var(--claude-text-secondary);
          
          p {
            margin: 0;
            font-size: 1rem;
          }
        }
      }
    }
    
    .chart-legend {
      padding: var(--claude-space-lg);
      padding-top: 0;
      display: flex;
      gap: var(--claude-space-xl);
      justify-content: center;
      
      .legend-item {
        display: flex;
        align-items: center;
        gap: var(--claude-space-sm);
        font-size: 0.9rem;
        color: var(--claude-text-secondary);
        
        .legend-color {
          width: 20px;
          height: 3px;
          border-radius: 2px;
          
          &.current {
            background: #3b82f6;
          }
          
          &.predicted {
            background: #ef4444;
            background-image: repeating-linear-gradient(
              90deg,
              #ef4444,
              #ef4444 5px,
              transparent 5px,
              transparent 10px
            );
          }
        }
      }
      
      @media (max-width: 480px) {
        flex-direction: column;
        align-items: center;
        gap: var(--claude-space);
      }
    }
    
    .fallback-table {
      padding: var(--claude-space-lg);
      border-top: 1px solid var(--claude-border);
      
      h4 {
        margin: 0 0 var(--claude-space-lg) 0;
        color: var(--claude-text-primary);
        font-size: 1rem;
      }
      
      .simple-data-list {
        display: flex;
        flex-direction: column;
        gap: var(--claude-space-sm);
        
        .data-item {
          display: grid;
          grid-template-columns: 1fr 1fr 1fr;
          gap: var(--claude-space);
          padding: var(--claude-space);
          background: var(--claude-bg-secondary);
          border-radius: var(--claude-radius);
          font-size: 0.9rem;
          
          .date {
            color: var(--claude-text-secondary);
            font-family: monospace;
          }
          
          .current-price {
            color: var(--claude-text-primary);
            font-family: monospace;
            font-weight: 600;
          }
          
          .predicted-price {
            color: var(--claude-primary);
            font-family: monospace;
            font-weight: 600;
          }
        }
      }
    }
  }
}

// 状态组件样式
.status-section {
  display: flex;
  justify-content: center;
  margin: var(--claude-space-xl) auto;
  
  .status-card {
    background: var(--claude-bg-primary);
    border: 1px solid var(--claude-border);
    border-radius: var(--claude-radius-lg);
    padding: var(--claude-space-xl);
    box-shadow: var(--claude-shadow);
    text-align: center;
    max-width: 500px;
    
    @media (max-width: 480px) {
      padding: var(--claude-space-lg);
    }
  }
  
  .status-text {
    font-size: 1rem;
    color: var(--claude-text-secondary);
    margin: var(--claude-space-lg) 0;
    line-height: 1.6;
  }
  
  &.loading {
    .loading-spinner {
      width: 48px;
      height: 48px;
      border: 3px solid var(--claude-bg-tertiary);
      border-top: 3px solid var(--claude-primary);
      border-radius: 50%;
      margin: 0 auto var(--claude-space-lg);
      animation: spin 1s linear infinite;
    }
  }
  
  &.error {
    .status-card {
      border-left: 4px solid var(--claude-danger);
    }
    
    .error-icon {
      font-size: 3rem;
      margin-bottom: var(--claude-space);
    }
  }
  
  &.empty {
    .empty-icon {
      font-size: 3rem;
      margin-bottom: var(--claude-space);
      opacity: 0.6;
    }
    
    .empty-description {
      color: var(--claude-text-tertiary);
      margin: var(--claude-space) 0;
    }
  }
}

.retry-button {
  background: var(--claude-primary);
  color: white;
  border: none;
  padding: var(--claude-space) var(--claude-space-lg);
  border-radius: var(--claude-radius);
  font-weight: 600;
  cursor: pointer;
  margin-top: var(--claude-space);
  
  &:hover {
    background: #2563eb;
  }
}

// 模态框样式
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: var(--claude-space-lg);
}

.modal-content {
  background: var(--claude-bg-primary);
  border-radius: var(--claude-radius-lg);
  box-shadow: var(--claude-shadow-lg);
  max-width: 600px;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--claude-space-lg);
  border-bottom: 1px solid var(--claude-border);
  
  .modal-title {
    font-size: 1.2rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    margin: 0;
  }
  
  .modal-close {
    background: none;
    border: none;
    font-size: 1.2rem;
    color: var(--claude-text-secondary);
    cursor: pointer;
    padding: var(--claude-space-sm);
    
    &:hover {
      color: var(--claude-text-primary);
    }
  }
}

.modal-body {
  padding: var(--claude-space-lg);
}

.indicator-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--claude-space-lg);
  
  .indicator-item {
    padding: var(--claude-space-lg);
    background: var(--claude-bg-tertiary);
    border-radius: var(--claude-radius);
    
    .indicator-label {
      font-size: 0.9rem;
      color: var(--claude-text-secondary);
      margin-bottom: var(--claude-space-sm);
    }
    
    .indicator-value {
      font-size: 1.1rem;
      font-weight: 600;
      color: var(--claude-primary);
      font-family: monospace;
    }
  }
}

// 动画
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>