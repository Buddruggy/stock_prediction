<template>
  <div class="home">
    <!-- 预测卡片 -->
    <div class="cards-grid">
      <div 
        v-for="(prediction, code) in displayData" 
        :key="code"
        class="card"
      >
        <div class="card-header">
          <h3 class="index-name">{{ prediction.name }}</h3>
          <span class="index-code">{{ code.toUpperCase() }}</span>
        </div>
        
        <div class="card-content">
          <div class="price-item">
            <span class="label">当前</span>
            <span class="value">{{ prediction.current?.toFixed(2) || '--' }}</span>
          </div>
          
          <div class="price-item">
            <span class="label">预测</span>
            <span class="value predicted">{{ prediction.predicted?.toFixed(2) || '--' }}</span>
          </div>
          
          <div class="price-item">
            <span class="label">涨跌</span>
            <span 
              class="value change" 
              :class="{ positive: prediction.change > 0, negative: prediction.change < 0 }"
            >
              {{ formatChange(prediction.change, prediction.changePercent) }}
            </span>
          </div>
          
          <div class="price-item">
            <span class="label">信心度</span>
            <span class="value confidence">{{ prediction.confidence?.toFixed(0) || '--' }}%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="status loading">
      <div class="spinner"></div>
      <span>正在获取数据...</span>
    </div>

    <!-- 使用mock数据提示 -->
    <div v-if="usingMockData" class="status info">
      <span>📊 当前显示模拟数据，实际数据加载中...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

// Mock数据
const mockData = {
  sh000001: {
    name: '上证指数',
    market: 'Shanghai',
    current: 3142.56,
    predicted: 3168.23,
    change: 25.67,
    changePercent: 0.82,
    confidence: 78.5
  },
  sz399001: {
    name: '深证成指',
    market: 'Shenzhen',
    current: 10234.78,
    predicted: 10187.45,
    change: -47.33,
    changePercent: -0.46,
    confidence: 72.3
  },
  sz399006: {
    name: '创业板指',
    market: 'ChiNext',
    current: 2156.89,
    predicted: 2178.12,
    change: 21.23,
    changePercent: 0.98,
    confidence: 65.8
  },
  sh000688: {
    name: '科创50',
    market: 'STAR50',
    current: 987.45,
    predicted: 994.67,
    change: 7.22,
    changePercent: 0.73,
    confidence: 69.2
  }
}

const predictions = ref({})
const loading = ref(false)
const usingMockData = ref(false)

// 显示的数据：优先使用真实数据，失败时使用mock数据
const displayData = computed(() => {
  return Object.keys(predictions.value).length > 0 ? predictions.value : mockData
})

// 格式化涨跌显示
const formatChange = (change, changePercent) => {
  if (change === undefined || changePercent === undefined) return '--'
  const sign = change > 0 ? '+' : ''
  return `${sign}${change.toFixed(2)} (${sign}${changePercent.toFixed(2)}%)`
}

const fetchPredictions = async () => {
  loading.value = true
  
  try {
    const response = await axios.get('/api/v1/predict/all', {
      timeout: 5000 // 5秒超时
    })
    
    if (response.data.code === 200) {
      predictions.value = response.data.data
      usingMockData.value = false
    } else {
      console.warn('API返回错误:', response.data.message)
      usingMockData.value = true
    }
  } catch (err) {
    console.warn('获取预测数据失败，使用模拟数据:', err.message)
    usingMockData.value = true
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPredictions()
  // 每30秒重试一次获取真实数据
  setInterval(fetchPredictions, 30 * 1000)
})
</script>

<style lang="scss" scoped>
.home {
  width: 100%;
}


// 卡片网格
.cards-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 1rem;
    max-width: 400px;
  }
  
  @media (max-width: 480px) {
    gap: 0.75rem;
    max-width: 100%;
  }
}

// 卡片样式 - 参考Claude官网风格
.card {
  background: #ffffff;
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid #e5e7eb;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  
  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
  
  @media (max-width: 768px) {
    padding: 1.25rem;
  }
  
  @media (max-width: 480px) {
    padding: 1rem;
    border-radius: 6px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
  
  .index-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0;
    letter-spacing: -0.025em;
    
    @media (max-width: 480px) {
      font-size: 1rem;
    }
  }
  
  .index-code {
    font-size: 0.75rem;
    color: #6b7280;
    background: #f9fafb;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-weight: 500;
    border: 1px solid #e5e7eb;
    
    @media (max-width: 480px) {
      font-size: 0.7rem;
      padding: 0.2rem 0.4rem;
    }
  }
}

.card-content {
  display: grid;
  gap: 0.75rem;
  
  @media (max-width: 480px) {
    gap: 0.5rem;
  }
}

.price-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  
  .label {
    font-size: 0.875rem;
    color: #6b7280;
    font-weight: 500;
    
    @media (max-width: 480px) {
      font-size: 0.8rem;
    }
  }
  
  .value {
    font-size: 1rem;
    font-weight: 600;
    color: #1a1a1a;
    
    &.predicted {
      color: #2563eb;
    }
    
    &.change {
      &.positive {
        color: #059669;
      }
      
      &.negative {
        color: #dc2626;
      }
    }
    
    &.confidence {
      color: #d97706;
    }
    
    @media (max-width: 480px) {
      font-size: 0.9rem;
    }
  }
}

// 状态提示
.status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1.5rem;
  margin: 2rem auto;
  max-width: 400px;
  border-radius: 8px;
  font-size: 0.9rem;
  
  &.loading {
    background: #f8f9fa;
    color: #6c757d;
  }
  
  &.info {
    background: #e7f3ff;
    color: #0066cc;
    border: 1px solid #b3d9ff;
  }
  
  @media (max-width: 768px) {
    margin: 1.5rem auto;
    padding: 1.25rem;
    font-size: 0.85rem;
  }
  
  @media (max-width: 480px) {
    margin: 1rem auto;
    padding: 1rem;
    font-size: 0.8rem;
  }
}

// 加载动画
.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #e9ecef;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  
  @media (max-width: 480px) {
    width: 16px;
    height: 16px;
  }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
