<template>
  <div class="claude-home">
    <!-- Claude 风格英雄区域 -->
    <div class="hero-section">
      <h2 class="hero-title">AI 驱动的智能股指预测</h2>
      <div v-if="stats" class="stats-section">
        <span class="stats-text">已累计预测 {{ stats.total_predictions }} 次，成功率 {{ stats.success_rate }}%</span>
      </div>
    </div>

    <!-- 预测卡片网格 -->
    <div class="predictions-section" v-if="Object.keys(predictions).length > 0">
      <div class="predictions-grid">
        <div 
          v-for="(prediction, code) in predictions" 
          :key="code"
          class="prediction-card"
          :class="{ 'high-confidence': prediction.confidence > 80 }"
        >
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="index-info">
              <div class="index-name-row">
                <h3 class="index-name">{{ prediction.name }}</h3>
                <span class="confidence-badge">{{ prediction.confidence?.toFixed(0) || '--' }}%</span>
              </div>
              <span class="index-code">{{ code.toUpperCase() }}</span>
            </div>
            <div class="trend-badge" :class="getTrendClass(prediction.change)">
              <span class="trend-value">{{ getTrendText(prediction.change) }}</span>
            </div>
          </div>
          
          <!-- 卡片内容 -->
          <div class="card-body">
            <div class="price-section">
              <div class="price-item current-price">
                <span class="price-label">当前价格</span>
                <span class="price-value">{{ prediction.current?.toFixed(2) || '--' }}</span>
              </div>
              
              <div class="price-item predicted-price">
                <span class="price-label">预测价格</span>
                <span class="price-value">{{ prediction.predicted?.toFixed(2) || '--' }}</span>
                <div class="change-info">
                  <span 
                    class="change-value" 
                    :class="{ positive: prediction.change > 0, negative: prediction.change < 0 }"
                  >
                    {{ formatChange(prediction.change, prediction.changePercent) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="status-section loading">
      <div class="status-card">
        <div class="loading-spinner"></div>
        <div class="status-text">
          <span v-if="Object.keys(predictions).length === 0">正在获取预测数据...</span>
          <span v-else>正在加载更多指数... ({{ Object.keys(predictions).length }}/4)</span>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-if="error && !loading" class="status-section error">
      <div class="status-card">
        <div class="error-icon">⚠️</div>
        <div class="status-text">{{ error }}</div>
        <button @click="fetchPredictions" class="retry-button">
          重新加载
        </button>
      </div>
    </div>

    <!-- 空数据状态 -->
    <div v-if="!loading && !error && Object.keys(predictions).length === 0" class="status-section empty">
      <div class="status-card">
        <div class="empty-icon">📊</div>
        <div class="status-text">暂无预测数据</div>
        <button @click="fetchPredictions" class="retry-button">
          刷新数据
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const predictions = ref({})
const loading = ref(false)
const error = ref('')
const stats = ref(null)

// 格式化涨跌显示
const formatChange = (change, changePercent) => {
  if (change === undefined || changePercent === undefined) return '--'
  const sign = change > 0 ? '+' : ''
  return `${sign}${change.toFixed(2)} (${sign}${changePercent.toFixed(2)}%)`
}

// 获取趋势类名
const getTrendClass = (change) => {
  if (change > 0) return 'positive'
  if (change < 0) return 'negative'
  return 'neutral'
}

// 获取趋势文字
const getTrendText = (change) => {
  if (change > 0) return '看涨'
  if (change < 0) return '看跌'
  return '持平'
}

// 获取置信度类名
const getConfidenceClass = (confidence) => {
  if (confidence >= 80) return 'high'
  if (confidence >= 60) return 'medium'
  return 'low'
}

// 支持的指数列表
const indices = [
  { code: 'sh000001', name: '上证综指' },
  { code: 'sz399001', name: '深证成指' },
  { code: 'sz399006', name: '创业板指' },
  { code: 'sh000688', name: '科创50' }
]

// 获取预测统计信息
const fetchPredictionStats = async () => {
  try {
    const response = await axios.get('/api/v1/prediction-stats')
    if (response.data.code === 200) {
      stats.value = response.data.data
    }
  } catch (err) {
    console.warn('获取预测统计信息失败:', err)
  }
}

const fetchPredictions = async () => {
  loading.value = true
  error.value = ''
  predictions.value = {} // 清空之前的预测结果
  
  let hasAnySuccess = false
  let allErrors = []
  
  // 逐个获取每个指数的预测
  for (const index of indices) {
    try {
      console.log(`正在获取 ${index.name}(${index.code}) 的预测数据...`)
      
      const response = await axios.get(`/api/v1/predict/${index.code}`, {
        timeout: 60000 // 60秒超时，给单个指数预测足够时间
      })
      
      if (response.data.code === 200) {
        // 成功获取预测，立即更新UI
        predictions.value[index.code] = response.data.data
        hasAnySuccess = true
        console.log(`${index.name} 预测获取成功`)
      } else {
        console.warn(`${index.name} 预测失败: ${response.data.message}`)
        allErrors.push(`${index.name}: ${response.data.message}`)
      }
    } catch (err) {
      let errorMsg = ''
      if (err.code === 'ECONNABORTED') {
        errorMsg = '请求超时'
      } else if (err.response) {
        errorMsg = `服务器错误(${err.response.status}): ${err.response.data?.message || err.message}`
      } else {
        errorMsg = `网络错误: ${err.message}`
      }
      
      console.warn(`${index.name} 预测失败: ${errorMsg}`)
      allErrors.push(`${index.name}: ${errorMsg}`)
    }
    
    // 在每次请求之间稍作停顿，避免服务器压力过大
    if (indices.indexOf(index) < indices.length - 1) {
      await new Promise(resolve => setTimeout(resolve, 100))
    }
  }
  
  // 处理最终结果
  if (!hasAnySuccess) {
    error.value = `所有指数预测失败:\n${allErrors.join('\n')}`
  } else if (allErrors.length > 0) {
    // 有部分成功，显示部分错误但不影响成功的结果
    console.warn('部分指数预测失败:', allErrors)
  }
  
  loading.value = false
}

onMounted(() => {
  fetchPredictionStats()
  fetchPredictions()
})
</script>

<style lang="scss" scoped>
@use '../assets/styles/modern.scss' as *;

.claude-home {
  max-width: 1000px;
  margin: 0 auto;
}

// 英雄区域
.hero-section {
  text-align: center;
  margin-bottom: var(--claude-space-xl); /* 缩小间距 */
  padding: var(--claude-space-lg) 0; /* 缩小内边距 */
  
  .hero-title {
    font-size: 2.5rem; /* 缩小标题字体 */
    font-weight: 700;
    background: linear-gradient(135deg, var(--claude-primary), var(--claude-primary-light));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: var(--claude-space);
    letter-spacing: -0.04em;
    line-height: 1.1;
    
    @media (max-width: 768px) {
      font-size: 2rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1.6rem;
    }
  }
  
  .stats-section {
    .stats-text {
      font-size: 1.1rem;
      color: var(--claude-text-secondary);
      font-weight: 500;
    }
  }
}

// 预测区域
.predictions-section {
  .predictions-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: var(--claude-space-lg);
    max-width: 1000px;
    margin: 0 auto;
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
      gap: var(--claude-space);
      max-width: 400px;
    }
  }
}

// 预测卡片
.prediction-card {
  @include claude-card;
  padding: var(--claude-space-lg);
  transition: var(--claude-transition);
  animation: claude-fade-in 0.6s ease-out;
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, var(--claude-primary), var(--claude-primary-light));
    opacity: 0;
    transition: var(--claude-transition);
  }
  
  &:hover {
    transform: translateY(-4px); /* 缩小悬停位移 */
    box-shadow: var(--claude-shadow-lg);
    
    &::before {
      opacity: 1;
    }
  }
  
  &.high-confidence {
    border-color: var(--claude-success);
    
    &::before {
      background: linear-gradient(90deg, var(--claude-success), var(--claude-accent));
      opacity: 1;
    }
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space-lg);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--claude-space);
  
  .index-info {
    flex: 1;
  }
  
  .index-name-row {
    display: flex;
    align-items: center;
    gap: var(--claude-space-sm);
    margin-bottom: var(--claude-space-sm);
  }
  
  .index-name {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    margin: 0;
    
    @media (max-width: 480px) {
      font-size: 1.25rem;
    }
  }
  
  .confidence-badge {
    background: var(--claude-primary);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: var(--claude-radius);
    font-size: 0.75rem;
    font-weight: 600;
    font-family: var(--claude-font-mono);
    
    @media (max-width: 480px) {
      font-size: 0.7rem;
      padding: 0.2rem 0.4rem;
    }
  }
  
  .index-code {
    display: inline-block;
    background: var(--claude-bg-tertiary);
    color: var(--claude-text-secondary);
    padding: 0.375rem 0.875rem;
    border-radius: var(--claude-radius-lg);
    font-size: 0.8rem;
    font-weight: 500;
    font-family: var(--claude-font-mono);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  
  .trend-badge {
    padding: 0.5rem 1rem;
    border-radius: var(--claude-radius-xl);
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    
    &.positive {
      background: rgba(5, 150, 105, 0.1);
      color: var(--claude-success);
    }
    
    &.negative {
      background: rgba(220, 38, 38, 0.1);
      color: var(--claude-danger);
    }
    
    &.neutral {
      background: var(--claude-bg-tertiary);
      color: var(--claude-text-tertiary);
    }
    
    @media (max-width: 480px) {
      font-size: 0.75rem;
      padding: 0.375rem 0.75rem;
    }
  }
}

.card-body {
  display: flex;
  flex-direction: column;
}

.price-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--claude-space);
  
  @media (max-width: 480px) {
    grid-template-columns: 1fr;
    gap: var(--claude-space-sm);
  }
  
  .price-item {
    text-align: center;
    padding: var(--claude-space);
    background: var(--claude-bg-tertiary);
    border-radius: var(--claude-radius-lg);
    
    .price-label {
      display: block;
      font-size: 0.9rem;
      color: var(--claude-text-secondary);
      margin-bottom: var(--claude-space-xs);
      font-weight: 500;
    }
    
    .price-value {
      display: block;
      font-size: 1.5rem;
      font-weight: 600;
      font-family: var(--claude-font-mono);
      color: var(--claude-text-primary);
      margin-bottom: var(--claude-space-xs);
      
      @media (max-width: 480px) {
        font-size: 1.25rem;
      }
    }
    
    &.predicted-price .price-value {
      color: var(--claude-primary);
    }
    
    .change-info {
      margin-top: var(--claude-space-xs);
      
      .change-value {
        font-size: 0.85rem;
        font-weight: 600;
        font-family: var(--claude-font-mono);
        
        &.positive {
          color: var(--claude-success);
        }
        
        &.negative {
          color: var(--claude-danger);
        }
        
        @media (max-width: 480px) {
          font-size: 0.8rem;
        }
      }
    }
    
    @media (max-width: 480px) {
      padding: var(--claude-space);
    }
  }
}



// 状态组件
.status-section {
  display: flex;
  justify-content: center;
  margin: var(--claude-space-xl) auto; /* 缩小间距 */
  
  .status-card {
    @include claude-card;
    text-align: center;
    padding: var(--claude-space-xl); /* 缩小内边距 */
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
    white-space: pre-line;
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
  }
}

.retry-button {
  @include claude-button-primary;
  margin-top: var(--claude-space);
}

// 动画
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

</style>