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
          <label class="filter-label">选择指数:</label>
          <select v-model="selectedIndex" @change="fetchHistoricalData" class="filter-select">
            <option value="all">全部指数</option>
            <option value="sh000001">上证综指</option>
            <option value="sz399001">深证成指</option>
            <option value="sz399006">创业板指</option>
            <option value="sh000688">科创50</option>
          </select>
        </div>
        
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

    <!-- 历史数据表格 -->
    <div v-else class="history-tables">
      <div v-for="(predictions, indexCode) in historicalData" :key="indexCode" class="table-section">
        <div class="table-header">
          <h3 class="table-title">{{ getIndexName(indexCode) }}</h3>
          <div class="table-stats">
            <span class="stat-item">记录数: {{ predictions.length }}</span>
            <span class="stat-item">平均置信度: {{ getAvgConfidence(predictions) }}%</span>
          </div>
        </div>
        
        <div class="table-container">
          <table class="history-table">
            <thead>
              <tr>
                <th>预测日期</th>
                <th>收盘价</th>
                <th>涨跌额</th>
                <th>涨跌幅</th>
                <th>置信度</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="prediction in predictions" :key="prediction.id || prediction.timestamp" class="table-row">
                <td class="date-cell">{{ formatDate(prediction.timestamp || prediction.prediction_date) }}</td>
                <td class="price-cell">{{ prediction.close_price?.toFixed(2) || '--' }}</td>
                <td class="change-cell" :class="getChangeClass(prediction.change)">
                  {{ formatChange(prediction.change) }}
                </td>
                <td class="percent-cell" :class="getChangeClass(prediction.change_percent)">
                  {{ formatPercent(prediction.change_percent) }}
                </td>
                <td class="confidence-cell">
                  <div class="confidence-bar-container">
                    <div 
                      class="confidence-bar"
                      :class="getConfidenceClass(prediction.confidence)"
                      :style="{ width: `${prediction.confidence || 0}%` }"
                    ></div>
                    <span class="confidence-text">{{ prediction.confidence?.toFixed(1) || '0' }}%</span>
                  </div>
                </td>
                <td>
                  <button 
                    @click="showIndicators(prediction)"
                    class="indicators-button"
                    v-if="prediction.technical_indicators"
                  >
                    查看技术指标
                  </button>
                  <span v-else class="text-muted">无数据</span>
                </td>
              </tr>
            </tbody>
          </table>
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
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

// 响应式数据
const loading = ref(false)
const error = ref('')
const historicalData = ref({})
const selectedIndex = ref('all')
const selectedDays = ref(30)
const showModal = ref(false)
const selectedPrediction = ref(null)

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
    let url = '/api/v1/predict/history'
    
    if (selectedIndex.value === 'all') {
      url += '/all'
    } else {
      url += `/${selectedIndex.value}`
    }
    
    url += `?days=${selectedDays.value}`
    
    const response = await axios.get(url)
    
    if (response.data.code === 200) {
      if (selectedIndex.value === 'all') {
        historicalData.value = response.data.data || {}
      } else {
        // 单个指数的数据需要包装成对象格式
        historicalData.value = {
          [selectedIndex.value]: response.data.data || []
        }
      }
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

.history-tables {
  .table-section {
    background: var(--claude-bg-primary);
    border: 1px solid var(--claude-border);
    border-radius: var(--claude-radius-lg);
    padding: 0;
    box-shadow: var(--claude-shadow);
    margin-bottom: var(--claude-space-xl);
    overflow: hidden;
    
    .table-header {
      padding: var(--claude-space-lg);
      border-bottom: 1px solid var(--claude-border);
      background: var(--claude-bg-tertiary);
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: var(--claude-space);
      
      .table-title {
        font-size: 1.2rem;
        font-weight: 600;
        color: var(--claude-text-primary);
        margin: 0;
      }
      
      .table-stats {
        display: flex;
        gap: var(--claude-space-lg);
        
        .stat-item {
          font-size: 0.9rem;
          color: var(--claude-text-secondary);
        }
      }
    }
    
    .table-container {
      overflow-x: auto;
    }
    
    .history-table {
      width: 100%;
      border-collapse: collapse;
      
      th, td {
        padding: var(--claude-space-lg);
        text-align: left;
        border-bottom: 1px solid var(--claude-border);
      }
      
      th {
        background: var(--claude-bg-tertiary);
        font-weight: 600;
        color: var(--claude-text-primary);
        font-size: 0.9rem;
        white-space: nowrap;
      }
      
      .table-row {
        &:hover {
          background: var(--claude-bg-tertiary);
        }
      }
      
      .date-cell {
        font-family: monospace;
        color: var(--claude-text-secondary);
      }
      
      .price-cell {
        font-family: monospace;
        font-weight: 600;
      }
      
      .change-cell, .percent-cell {
        font-family: monospace;
        font-weight: 600;
        
        &.positive {
          color: var(--claude-success);
        }
        
        &.negative {
          color: var(--claude-danger);
        }
      }
      
      .confidence-cell {
        .confidence-bar-container {
          display: flex;
          align-items: center;
          gap: var(--claude-space);
          
          .confidence-bar {
            height: 8px;
            border-radius: 4px;
            min-width: 40px;
            max-width: 80px;
            
            &.high {
              background: var(--claude-success);
            }
            
            &.medium {
              background: var(--claude-warning);
            }
            
            &.low {
              background: var(--claude-danger);
            }
          }
          
          .confidence-text {
            font-family: monospace;
            font-size: 0.85rem;
            font-weight: 600;
          }
        }
      }
      
      .indicators-button {
        background: var(--claude-bg-secondary);
        color: var(--claude-text-primary);
        border: 1px solid var(--claude-border);
        padding: var(--claude-space-sm) var(--claude-space);
        border-radius: var(--claude-radius);
        font-size: 0.8rem;
        cursor: pointer;
        
        &:hover {
          background: var(--claude-bg-tertiary);
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