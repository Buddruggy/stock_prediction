<template>
  <div class="about">
    <div class="about-content">
      <h2>关于智投预测</h2>
      
      <div class="info-cards">
        <div class="info-card">
          <el-icon><TrendCharts /></el-icon>
          <h3>AI技术</h3>
          <p>采用先进的机器学习算法，结合多种技术指标进行智能预测</p>
        </div>
        
        <div class="info-card">
          <el-icon><DataAnalysis /></el-icon>
          <h3>实时数据</h3>
          <p>获取最新的股票市场数据，确保预测结果的时效性</p>
        </div>
        
        <div class="info-card">
          <el-icon><Monitor /></el-icon>
          <h3>专业服务</h3>
          <p>专注于中国主要股票指数，提供专业的预测分析服务</p>
        </div>
      </div>
      
      <div class="disclaimer-card">
        <el-icon><WarningFilled /></el-icon>
        <h3>风险提示</h3>
        <p>本平台基于AI技术提供指数预测，结果仅供参考，不构成投资建议。股市有风险，投资需谨慎。</p>
      </div>
      
      <div class="tech-stack">
        <h3>技术栈</h3>
        <div class="tech-tags">
          <span class="tech-tag">Vue.js 3</span>
          <span class="tech-tag">FastAPI</span>
          <span class="tech-tag">Python</span>
          <span class="tech-tag">Machine Learning</span>
          <span class="tech-tag">Docker</span>
          <span class="tech-tag">Nginx</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { TrendCharts, DataAnalysis, Monitor, WarningFilled } from '@element-plus/icons-vue'
</script>

<style lang="scss" scoped>
.about {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.about-content {
  h2 {
    color: white;
    font-size: 2.5rem;
    text-align: center;
    margin-bottom: 40px;
  }
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.info-card {
  background: rgba(255,255,255,0.1);
  border-radius: 16px;
  padding: 30px;
  text-align: center;
  backdrop-filter: blur(10px);
  transition: transform 0.3s ease;
  
  &:hover {
    transform: translateY(-4px);
  }
  
  .el-icon {
    font-size: 48px;
    color: #409eff;
    margin-bottom: 16px;
  }
  
  h3 {
    color: white;
    font-size: 1.25rem;
    margin-bottom: 12px;
  }
  
  p {
    color: rgba(255,255,255,0.8);
    line-height: 1.6;
    margin: 0;
  }
}

.disclaimer-card {
  background: rgba(255, 193, 7, 0.1);
  border: 1px solid rgba(255, 193, 7, 0.3);
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  margin-bottom: 40px;
  
  .el-icon {
    font-size: 32px;
    color: #ffc107;
    margin-bottom: 12px;
  }
  
  h3 {
    color: #ffc107;
    font-size: 1.125rem;
    margin-bottom: 12px;
  }
  
  p {
    color: rgba(255,255,255,0.9);
    line-height: 1.6;
    margin: 0;
  }
}

.tech-stack {
  background: rgba(255,255,255,0.05);
  border-radius: 16px;
  padding: 30px;
  text-align: center;
  
  h3 {
    color: white;
    font-size: 1.25rem;
    margin-bottom: 20px;
  }
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.tech-tag {
  background: rgba(64, 158, 255, 0.2);
  color: #409eff;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 500;
  border: 1px solid rgba(64, 158, 255, 0.3);
}

@media (max-width: 768px) {
  .about-content h2 {
    font-size: 2rem;
  }
  
  .info-cards {
    grid-template-columns: 1fr;
  }
  
  .tech-tags {
    justify-content: flex-start;
  }
}
</style>
