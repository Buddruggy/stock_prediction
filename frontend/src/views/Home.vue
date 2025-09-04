<template>
  <div class="home">
    <div class="hero-section">
      <h2>AI股市指数预测</h2>
      <p>基于机器学习技术，为您提供专业的股票指数预测服务</p>
    </div>

    <div class="predictions-grid">
      <div 
        v-for="(prediction, code) in predictions" 
        :key="code"
        class="prediction-card"
        :class="`prediction-${code}`"
      >
        <div class="card-header">
          <h3>{{ prediction.name }}</h3>
          <span class="market">{{ prediction.market }}</span>
        </div>
        
        <div class="card-body">
          <div class="current-price">
            <label>当前价格</label>
            <span class="price">{{ prediction.current?.toFixed(2) || '--' }}</span>
          </div>
          
          <div class="predicted-price">
            <label>明日预测</label>
            <span class="price predicted">{{ prediction.predicted?.toFixed(2) || '--' }}</span>
          </div>
          
          <div class="change" :class="{ positive: prediction.change > 0, negative: prediction.change < 0 }">
            <label>预测涨跌</label>
            <span>{{ prediction.change > 0 ? '+' : '' }}{{ prediction.change?.toFixed(2) || '--' }} ({{ prediction.changePercent?.toFixed(2) || '--' }}%)</span>
          </div>
          
          <div class="confidence">
            <label>AI信心指数</label>
            <span>{{ prediction.confidence?.toFixed(1) || '--' }}%</span>
          </div>
        </div>
      </div>
    </div>

    <div class="loading" v-if="loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>正在获取预测数据...</span>
    </div>

    <div class="error" v-if="error">
      <el-icon><Warning /></el-icon>
      <span>{{ error }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Loading, Warning } from '@element-plus/icons-vue'
import axios from 'axios'

const predictions = ref({})
const loading = ref(false)
const error = ref('')

const fetchPredictions = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await axios.get('/api/v1/predict/all')
    if (response.data.code === 200) {
      predictions.value = response.data.data
    } else {
      error.value = response.data.message || '获取数据失败'
    }
  } catch (err) {
    error.value = '网络连接失败，请稍后重试'
    console.error('获取预测数据失败:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPredictions()
  // 每5分钟刷新一次数据
  setInterval(fetchPredictions, 5 * 60 * 1000)
})
</script>

<style lang="scss" scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.hero-section {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 20px;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  
  h2 {
    font-size: 2.5rem;
    color: white;
    margin-bottom: 16px;
    font-weight: 700;
  }
  
  p {
    font-size: 1.2rem;
    color: rgba(255,255,255,0.8);
    margin: 0;
  }
}

.predictions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.prediction-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border-left: 4px solid #409eff;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 30px rgba(0,0,0,0.15);
  }
  
  &.prediction-sh000001 {
    border-left-color: #e74c3c;
  }
  
  &.prediction-sz399001 {
    border-left-color: #2ecc71;
  }
  
  &.prediction-sz399006 {
    border-left-color: #f39c12;
  }
  
  &.prediction-sh000688 {
    border-left-color: #9b59b6;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  
  h3 {
    margin: 0;
    font-size: 1.25rem;
    color: #2c3e50;
    font-weight: 600;
  }
  
  .market {
    font-size: 0.875rem;
    color: #7f8c8d;
    background: #ecf0f1;
    padding: 4px 8px;
    border-radius: 4px;
  }
}

.card-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  
  > div {
    display: flex;
    flex-direction: column;
    
    label {
      font-size: 0.875rem;
      color: #7f8c8d;
      margin-bottom: 4px;
    }
    
    span {
      font-size: 1.125rem;
      font-weight: 600;
      color: #2c3e50;
    }
  }
  
  .predicted span {
    color: #409eff;
  }
  
  .change {
    &.positive span {
      color: #67c23a;
    }
    
    &.negative span {
      color: #f56c6c;
    }
  }
  
  .confidence span {
    color: #e6a23c;
  }
}

.loading, .error {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
  text-align: center;
  color: white;
  
  .el-icon {
    font-size: 24px;
  }
}

.error {
  color: #f56c6c;
}

@media (max-width: 768px) {
  .hero-section h2 {
    font-size: 2rem;
  }
  
  .predictions-grid {
    grid-template-columns: 1fr;
  }
  
  .card-body {
    grid-template-columns: 1fr;
  }
}
</style>
