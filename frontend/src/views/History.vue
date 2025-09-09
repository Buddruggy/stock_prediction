<template>
  <div class="modern-history">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">历史预测记录</h2>
      <p class="page-subtitle">查看各指数的历史预测数据和准确性分析</p>
    </div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <div class="filter-section">
        <div class="filter-item">
          <label class="filter-label">选择指数：</label>
          <select v-model="selectedIndex" @change="fetchHistoricalData" class="filter-select">
            <option value="all">全部指数</option>
            <option value="sh000001">上证综指</option>
            <option value="sz399001">深证成指</option>
            <option value="sz399006">创业板指</option>
            <option value="sh000688">科创50</option>
          </select>
        </div>
        
        <div class="filter-item">
          <label class="filter-label">时间范围：</label>
          <select v-model="selectedDays" @change="fetchHistoricalData" class="filter-select">
            <option :value="7">最近7天</option>
            <option :value="15">最近15天</option>
            <option :value="30">最近30天</option>
            <option :value="60">最近60天</option>
            <option :value="90">最近90天</option>
          </select>
        </div>
        
        <button @click="fetchHistoricalData" class="refresh-button" :disabled="loading">
          <span v-if="loading">🔄</span>
          <span v-else>📊</span>
          {{ loading ? '加载中...' : '刷新数据' }}
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="status-section loading">
      <div class="status-card">
        <div class="loading-spinner"></div>
        <div class="status-text">正在获取历史预测数据...</div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-if="error && !loading" class="status-section error">
      <div class="status-card">
        <div class="error-icon">⚠️</div>
        <div class="status-text">{{ error }}</div>
        <button @click="fetchHistoricalData" class="retry-button">
          重新加载
        </button>
      </div>
    </div>

    <!-- 数据展示 -->
    <div v-if="!loading && !error && hasData" class="data-section">
      <!-- 统计概览 -->
      <div class="stats-overview">
        <div class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-content">
            <div class="stat-number">{{ totalPredictions }}</div>
            <div class="stat-label">总预测次数</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">📈</div>
          <div class="stat-content">
            <div class="stat-number">{{ avgConfidence }}%</div>
            <div class="stat-label">平均置信度</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">📅</div>
          <div class="stat-content">
            <div class="stat-number">{{ selectedDays }}</div>
            <div class="stat-label">天数范围</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">🎯</div>
          <div class="stat-content">
            <div class="stat-number">{{ Object.keys(historicalData).length }}</div>
            <div class="stat-label">涉及指数</div>
          </div>
        </div>
      </div>

      <!-- 历史数据表格 -->
      <div class="history-tables">
        <div v-for="(predictions, indexCode) in historicalData" :key="indexCode" class="table-section">
          <div class="table-header">
            <h3 class="table-title">{{ getIndexName(indexCode) }}</h3>
            <div class="table-stats">
              <span class="stat-item">共 {{ predictions.length }} 条记录</span>
              <span class="stat-item">平均置信度: {{ getAvgConfidence(predictions) }}%</span>
            </div>
          </div>
          
          <div class="table-container">
            <table class="history-table">
              <thead>
                <tr>
                  <th>预测日期</th>
                  <th>当前价格</th>
                  <th>预测价格</th>
                  <th>预测涨跌</th>
                  <th>预测涨跌幅</th>
                  <th>置信度</th>
                  <th>技术指标</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(prediction, index) in predictions" :key="index" class="table-row">
                  <td class="date-cell">{{ formatDate(prediction.timestamp) }}</td>
                  <td class="price-cell">{{ prediction.current?.toFixed(2) || '--' }}</td>
                  <td class="price-cell">{{ prediction.predicted?.toFixed(2) || '--' }}</td>
                  <td class="change-cell" :class="getChangeClass(prediction.change)">
                    {{ formatChange(prediction.change) }}
                  </td>
                  <td class="percent-cell" :class="getChangeClass(prediction.change)">
                    {{ formatPercent(prediction.changePercent) }}
                  </td>
                  <td class="confidence-cell">
                    <div class="confidence-bar-container">
                      <div 
                        class="confidence-bar" 
                        :style="{ width: prediction.confidence + '%' }"
                        :class="getConfidenceClass(prediction.confidence)"
                      ></div>
                      <span class="confidence-text">{{ prediction.confidence?.toFixed(1) || '--' }}%</span>
                    </div>
                  </td>
                  <td class="indicators-cell">
                    <button @click="showIndicators(prediction)" class="indicators-button">
                      查看详情
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- 空数据状态 -->
    <div v-if="!loading && !error && !hasData" class="status-section empty">
      <div class="status-card">
        <div class="empty-icon">📊</div>
        <div class="status-text">暂无历史预测数据</div>
        <p class="empty-description">请等待系统生成预测数据，或调整筛选条件</p>
        <button @click="fetchHistoricalData" class="retry-button">
          刷新数据
        </button>
      </div>
    </div>

    <!-- 技术指标详情模态框 -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">技术指标详情</h3>
          <button @click="closeModal" class="modal-close">✕</button>
        </div>
        
        <div class="modal-body" v-if="selectedPrediction">
          <div class="indicator-grid">
            <div class="indicator-item">
              <div class="indicator-label">5日移动平均线 (MA5)</div>
              <div class="indicator-value">{{ selectedPrediction.technical_indicators?.ma_5?.toFixed(2) || '--' }}</div>
            </div>
            
            <div class="indicator-item">
              <div class="indicator-label">20日移动平均线 (MA20)</div>
              <div class="indicator-value">{{ selectedPrediction.technical_indicators?.ma_20?.toFixed(2) || '--' }}</div>
            </div>
            
            <div class="indicator-item">
              <div class="indicator-label">相对强弱指数 (RSI)</div>
              <div class="indicator-value">{{ selectedPrediction.technical_indicators?.rsi?.toFixed(2) || '--' }}</div>
            </div>
            
            <div class="indicator-item">
              <div class="indicator-label">波动率</div>
              <div class="indicator-value">{{ selectedPrediction.technical_indicators?.volatility?.toFixed(2) || '--' }}%</div>
            </div>
            
            <div class="indicator-item">
              <div class="indicator-label">趋势指标</div>
              <div class="indicator-value">{{ selectedPrediction.technical_indicators?.trend?.toFixed(2) || '--' }}%</div>
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
                  <p class="feature-desc">直观比较不同指数的历史表现</p>
                </div>
              </div>
              
              <div class="feature-card">
                <div class="feature-icon">🎯</div>
                <div class="feature-content">
                  <h5 class="feature-name">预测误差分析</h5>
                  <p class="feature-desc">详细分析预测误差的分布和趋势</p>
                </div>
              </div>
              
              <div class="feature-card">
                <div class="feature-icon">⚡</div>
                <div class="feature-content">
                  <h5 class="feature-name">模型性能评估</h5>
                  <p class="feature-desc">全面评估AI模型的预测性能</p>
                </div>
              </div>
            </div>
          </div>
          
          <div class="progress-section">
            <div class="progress-header">
              <span class="progress-label">开发进度</span>
              <span class="progress-value">75%</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill" style="width: 75%"></div>
            </div>
          </div>
          
          <div class="notice-footer">
            <div class="status-badge">
              <span class="status-dot"></span>
              <span class="status-text">开发中</span>
            </div>
            <span class="eta-text">预计完成时间：2024年Q2</span>
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
const getBarHeight = (index) => {
  const heights = ['60%', '80%', '45%', '90%', '70%', '55%']
  return heights[index - 1] || '50%'
}
</script>

<style lang="scss" scoped>
@import '../assets/styles/main.scss';

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
  @include claude-card;
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
    @include claude-button-primary;
    display: flex;
    align-items: center;
    gap: var(--claude-space-sm);
    
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
    @include claude-card;
    display: flex;
    align-items: center;
    gap: var(--claude-space-lg);
    padding: var(--claude-space-lg);
    
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
    @include claude-card;
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
        @include claude-button-secondary;
        font-size: 0.8rem;
        padding: var(--claude-space-sm) var(--claude-space);
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
    @include claude-card;
    text-align: center;
    padding: var(--claude-space-xl);
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
  @include claude-button-primary;
  margin-top: var(--claude-space);
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