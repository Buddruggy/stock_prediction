<template>
  <div class="history">
    <h2>历史数据</h2>
    <p>查看各指数的历史走势和预测准确性</p>
    
    <div class="coming-soon">
      <el-icon><DataAnalysis /></el-icon>
      <h3>功能开发中</h3>
      <p>历史数据分析功能正在开发中，敬请期待...</p>
    </div>
  </div>
</template>

<script setup>
import { DataAnalysis } from '@element-plus/icons-vue'
</script>

<style lang="scss" scoped>
.history {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  
  h2 {
    color: white;
    font-size: 2rem;
    margin-bottom: 16px;
  }
  
  > p {
    color: rgba(255,255,255,0.8);
    font-size: 1.1rem;
    margin-bottom: 60px;
  }
}

.coming-soon {
  background: rgba(255,255,255,0.1);
  border-radius: 16px;
  padding: 60px 40px;
  backdrop-filter: blur(10px);
  
  .el-icon {
    font-size: 64px;
    color: #409eff;
    margin-bottom: 20px;
  }
  
  h3 {
    color: white;
    font-size: 1.5rem;
    margin-bottom: 16px;
  }
  
  p {
    color: rgba(255,255,255,0.7);
    font-size: 1rem;
    margin: 0;
  }
}
</style>
