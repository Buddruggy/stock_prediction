<template>
  <div class="home">
    <!-- 预测卡片 -->
    <div class="cards-grid" v-if="Object.keys(predictions).length > 0">
      <div 
        v-for="(prediction, code) in predictions" 
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
            <span class="label">预测涨跌</span>
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

    <!-- 错误状态 -->
    <div v-if="error && !loading" class="status error">
      <span>⚠️ {{ error }}</span>
      <button @click="fetchPredictions" class="retry-btn">重试</button>
    </div>

    <!-- 空数据状态 -->
    <div v-if="!loading && !error && Object.keys(predictions).length === 0" class="status empty">
      <span>📊 暂无预测数据</span>
      <button @click="fetchPredictions" class="retry-btn">刷新</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const predictions = ref({})
const loading = ref(false)
const error = ref('')

// 格式化涨跌显示
const formatChange = (change, changePercent) => {
  if (change === undefined || changePercent === undefined) return '--'
  const sign = change > 0 ? '+' : ''
  return `${sign}${change.toFixed(2)} (${sign}${changePercent.toFixed(2)}%)`
}

const fetchPredictions = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await axios.get('/api/v1/predict/all', {
      timeout: 60000 // 60秒超时，与DeepSeek API超时保持一致
    })
    
    if (response.data.code === 200) {
      predictions.value = response.data.data
    } else {
      error.value = `API错误: ${response.data.message}`
    }
  } catch (err) {
    if (err.code === 'ECONNABORTED') {
      error.value = '请求超时：AI预测服务响应耗时较长，请检查网络或稍后重试'
    } else if (err.response) {
      error.value = `服务器错误: ${err.response.status} - ${err.response.data?.message || err.message}`
    } else {
      error.value = `网络错误: ${err.message}`
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPredictions()
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
  max-width: 500px;
  border-radius: 8px;
  font-size: 0.9rem;
  
  &.loading {
    background: #f8f9fa;
    color: #6c757d;
  }
  
  &.error {
    background: #fef2f2;
    color: #dc2626;
    border: 1px solid #fecaca;
    flex-direction: column;
    gap: 1rem;
  }
  
  &.empty {
    background: #f9fafb;
    color: #6b7280;
    border: 1px solid #e5e7eb;
    flex-direction: column;
    gap: 1rem;
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

// 重试按钮
.retry-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
  
  &:hover {
    background: #1d4ed8;
  }
  
  &:active {
    background: #1e40af;
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
