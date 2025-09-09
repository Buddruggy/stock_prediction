<template>
  <div class="modern-history">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">历史数据分析</h2>
      <p class="page-subtitle">查看各指数的历史走势和预测准确性统计</p>
    </div>
    
    <!-- 开发中状态 -->
    <div class="development-notice">
      <div class="notice-content">
        <div class="notice-icon">
          <div class="chart-placeholder">
            <div class="chart-bars">
              <div class="bar" v-for="i in 6" :key="i" :style="{ height: getBarHeight(i) }"></div>
            </div>
          </div>
        </div>
        
        <div class="notice-body">
          <h3 class="notice-title">功能开发中</h3>
          <p class="notice-description">
            历史数据分析功能正在紧张开发中，即将为您提供更全面的数据洞察
          </p>
          
          <div class="features-preview">
            <h4 class="features-title">即将推出的功能：</h4>
            <div class="features-grid">
              <div class="feature-card">
                <div class="feature-icon">📊</div>
                <div class="feature-content">
                  <h5 class="feature-name">预测准确率统计</h5>
                  <p class="feature-desc">实时统计各指数预测的准确率</p>
                </div>
              </div>
              
              <div class="feature-card">
                <div class="feature-icon">📈</div>
                <div class="feature-content">
                  <h5 class="feature-name">历史走势对比</h5>
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
// 生成图表条形高度
const getBarHeight = (index) => {
  const heights = ['60%', '80%', '45%', '90%', '70%', '55%']
  return heights[index - 1] || '50%'
}
</script>

<style lang="scss" scoped>
@use '../assets/styles/modern.scss' as *;

.modern-history {
  max-width: 1000px;
  margin: 0 auto;
}

// 页面标题
.page-header {
  text-align: center;
  margin-bottom: var(--claude-space-2xl);
  
  .page-title {
    font-size: 2.25rem;
    font-weight: 700;
    @include claude-text-gradient(var(--claude-secondary), var(--claude-primary));
    margin-bottom: var(--claude-space);
    letter-spacing: -0.025em;
    
    @media (max-width: 768px) {
      font-size: 1.875rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1.5rem;
    }
  }
  
  .page-subtitle {
    font-size: 1.125rem;
    color: var(--claude-text-secondary);
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
    
    @media (max-width: 768px) {
      font-size: 1rem;
    }
    
    @media (max-width: 480px) {
      font-size: 0.9rem;
      padding: 0 var(--claude-space-sm);
    }
  }
}

// 开发通知卡片
.development-notice {
  @include claude-card;
  padding: var(--claude-space-2xl);
  border-left: 4px solid var(--claude-warning);
  animation: claude-slide-up 0.6s ease-out;
  
  @media (max-width: 768px) {
    padding: var(--claude-space-xl);
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space-lg);
  }
}

.notice-content {
  display: flex;
  gap: var(--claude-space-xl);
  align-items: flex-start;
  
  @media (max-width: 768px) {
    flex-direction: column;
    gap: var(--claude-space-lg);
    text-align: center;
  }
}

.notice-icon {
  flex-shrink: 0;
  
  .chart-placeholder {
    width: 120px;
    height: 80px;
    background: var(--claude-bg-tertiary);
    border-radius: var(--claude-radius-lg);
    display: flex;
    align-items: flex-end;
    justify-content: center;
    padding: var(--claude-space);
    border: 2px solid var(--claude-border);
    
    .chart-bars {
      display: flex;
      align-items: flex-end;
      gap: 4px;
      height: 100%;
      
      .bar {
        width: 8px;
        background: linear-gradient(to top, var(--claude-primary), var(--claude-accent));
        border-radius: 2px;
        animation: bar-grow 1.5s ease-out;
        animation-delay: calc(var(--index, 0) * 0.1s);
        
        @for $i from 1 through 6 {
          &:nth-child(#{$i}) {
            --index: #{$i - 1};
          }
        }
      }
    }
    
    @media (max-width: 768px) {
      width: 100px;
      height: 60px;
    }
    
    @media (max-width: 480px) {
      width: 80px;
      height: 50px;
      margin: 0 auto;
    }
  }
}

.notice-body {
  flex: 1;
  
  .notice-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    margin-bottom: var(--claude-space);
    
    @media (max-width: 480px) {
      font-size: 1.25rem;
    }
  }
  
  .notice-description {
    font-size: 1rem;
    color: var(--claude-text-secondary);
    line-height: 1.6;
    margin-bottom: var(--claude-space-xl);
    
    @media (max-width: 480px) {
      font-size: 0.9rem;
    }
  }
}

// 功能预览
.features-preview {
  margin-bottom: var(--claude-space-xl);
  
  .features-title {
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    margin-bottom: var(--claude-space-md);
    
    @media (max-width: 480px) {
      font-size: 1rem;
    }
  }
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: var(--claude-space);
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }
}

.feature-card {
  background: var(--claude-bg-tertiary);
  border: 1px solid var(--claude-border-light);
  border-radius: var(--claude-radius-md);
  padding: var(--claude-space);
  display: flex;
  gap: var(--claude-space);
  align-items: flex-start;
  transition: var(--claude-transition);
  
  &:hover {
    background: var(--claude-bg-card);
    border-color: var(--claude-border);
    transform: translateY(-2px);
    box-shadow: var(--claude-shadow);
  }
  
  .feature-icon {
    font-size: 1.5rem;
    flex-shrink: 0;
    
    @media (max-width: 480px) {
      font-size: 1.25rem;
    }
  }
  
  .feature-content {
    .feature-name {
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--claude-text-primary);
      margin-bottom: var(--claude-space-xs);
      
      @media (max-width: 480px) {
        font-size: 0.8rem;
      }
    }
    
    .feature-desc {
      font-size: 0.75rem;
      color: var(--claude-text-secondary);
      line-height: 1.4;
      margin: 0;
      
      @media (max-width: 480px) {
        font-size: 0.7rem;
      }
    }
  }
}

// 进度条
.progress-section {
  margin-bottom: var(--claude-space-xl);
  
  .progress-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--claude-space-sm);
    
    .progress-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--claude-text-secondary);
    }
    
    .progress-value {
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--claude-warning);
      font-family: var(--claude-font-mono);
    }
  }
  
  .progress-bar {
    height: 8px;
    background: var(--claude-bg-tertiary);
    border-radius: var(--claude-radius);
    overflow: hidden;
    
    .progress-fill {
      height: 100%;
      background: linear-gradient(90deg, var(--claude-warning), #fbbf24);
      border-radius: var(--claude-radius);
      transition: width 1s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      
      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
        animation: progress-shine 2s infinite;
      }
    }
  }
}

// 通知页脚
.notice-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  @media (max-width: 480px) {
    flex-direction: column;
    gap: var(--claude-space-sm);
  }
  
  .status-badge {
    display: flex;
    align-items: center;
    gap: var(--claude-space-xs);
    background: rgba(245, 158, 11, 0.1);
    color: var(--claude-warning);
    padding: 0.375rem 0.75rem;
    border-radius: var(--claude-radius-lg);
    font-size: 0.75rem;
    font-weight: 500;
    
    .status-dot {
      width: 6px;
      height: 6px;
      background: var(--claude-warning);
      border-radius: 50%;
      animation: pulse-dot 2s infinite;
    }
  }
  
  .eta-text {
    font-size: 0.75rem;
    color: var(--claude-text-tertiary);
    font-style: italic;
  }
}

// 动画
@keyframes bar-grow {
  from {
    height: 0;
  }
  to {
    height: var(--final-height, 50%);
  }
}

@keyframes progress-shine {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

@keyframes pulse-dot {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.2);
  }
}
</style>