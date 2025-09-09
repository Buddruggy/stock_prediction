<template>
  <div class="modern-about">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        关于<span class="about-logo">
          <div class="mini-dog">
            <div class="mini-face"></div>
            <div class="mini-ear-left"></div>
            <div class="mini-ear-right"></div>
          </div>
          <span class="logo-text">GoGoTou</span>
        </span>
      </h2>
      <p class="page-subtitle">基于人工智能技术的股票指数预测平台</p>
    </div>
    
    <!-- 功能特性卡片 -->
    <div class="features-section">
      <div class="features-grid">
        <div class="feature-card ai-feature">
          <div class="feature-header">
            <div class="feature-icon">
              <div class="ai-brain">
                <div class="brain-core"></div>
                <div class="brain-waves">
                  <div class="wave" v-for="i in 3" :key="i"></div>
                </div>
              </div>
            </div>
            <h3 class="feature-title">AI 技术</h3>
          </div>
          <div class="feature-content">
            <p class="feature-description">
              采用先进的机器学习算法，结合多种技术指标进行智能预测分析，
              为投资决策提供科学依据。
            </p>
            <div class="tech-highlights">
              <span class="highlight-tag">随机森林</span>
              <span class="highlight-tag">技术指标</span>
              <span class="highlight-tag">深度学习</span>
            </div>
          </div>
        </div>
        
        <div class="feature-card data-feature">
          <div class="feature-header">
            <div class="feature-icon">
              <div class="data-chart">
                <div class="chart-line"></div>
                <div class="data-points">
                  <div class="point" v-for="i in 4" :key="i"></div>
                </div>
              </div>
            </div>
            <h3 class="feature-title">实时数据</h3>
          </div>
          <div class="feature-content">
            <p class="feature-description">
              获取最新的股票市场数据，确保预测结果的时效性和准确性，
              支持多个主要指数的实时监控。
            </p>
            <div class="tech-highlights">
              <span class="highlight-tag">实时API</span>
              <span class="highlight-tag">数据清洗</span>
              <span class="highlight-tag">质量控制</span>
            </div>
          </div>
        </div>
        
        <div class="feature-card service-feature">
          <div class="feature-header">
            <div class="feature-icon">
              <div class="service-monitor">
                <div class="monitor-screen">
                  <div class="screen-content">
                    <div class="screen-line" v-for="i in 3" :key="i"></div>
                  </div>
                </div>
              </div>
            </div>
            <h3 class="feature-title">专业服务</h3>
          </div>
          <div class="feature-content">
            <p class="feature-description">
              专注于中国主要股票指数，提供专业的预测分析服务，
              帮助投资者做出更明智的投资决策。
            </p>
            <div class="tech-highlights">
              <span class="highlight-tag">中国市场</span>
              <span class="highlight-tag">指数分析</span>
              <span class="highlight-tag">投资助手</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 风险提示 -->
    <div class="risk-warning">
      <div class="warning-header">
        <div class="warning-icon">
          <div class="warning-symbol">⚠️</div>
        </div>
        <h3 class="warning-title">重要风险提示</h3>
      </div>
      <div class="warning-content">
        <p class="warning-text">
          本平台基于AI技术提供指数预测，结果仅供参考，不构成投资建议。
          <strong>股市有风险，投资需谨慎。</strong>
        </p>
        <div class="warning-points">
          <div class="warning-point">
            <span class="point-icon">•</span>
            <span class="point-text">预测结果基于历史数据，无法保证未来准确性</span>
          </div>
          <div class="warning-point">
            <span class="point-icon">•</span>
            <span class="point-text">投资决策应结合多方面因素综合考虑</span>
          </div>
          <div class="warning-point">
            <span class="point-icon">•</span>
            <span class="point-text">请根据自身风险承受能力合理投资</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 技术栈展示 -->
    <div class="tech-stack-section">
      <h3 class="section-title">技术架构</h3>
      <div class="tech-categories">
        <div class="tech-category frontend-tech">
          <div class="category-header">
            <div class="category-icon">🎨</div>
            <h4 class="category-title">前端技术</h4>
          </div>
          <div class="tech-list">
            <div class="tech-item">
              <span class="tech-name">Vue.js 3</span>
              <span class="tech-desc">现代化前端框架</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">Vite</span>
              <span class="tech-desc">高性能构建工具</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">响应式设计</span>
              <span class="tech-desc">适配多端设备</span>
            </div>
          </div>
        </div>
        
        <div class="tech-category backend-tech">
          <div class="category-header">
            <div class="category-icon">⚙️</div>
            <h4 class="category-title">后端技术</h4>
          </div>
          <div class="tech-list">
            <div class="tech-item">
              <span class="tech-name">Go Language</span>
              <span class="tech-desc">高性能后端服务</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">Gin Framework</span>
              <span class="tech-desc">轻量级Web框架</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">RESTful API</span>
              <span class="tech-desc">标准化接口设计</span>
            </div>
          </div>
        </div>
        
        <div class="tech-category ai-tech">
          <div class="category-header">
            <div class="category-icon">🤖</div>
            <h4 class="category-title">AI & 机器学习</h4>
          </div>
          <div class="tech-list">
            <div class="tech-item">
              <span class="tech-name">随机森林</span>
              <span class="tech-desc">主要预测算法</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">技术指标</span>
              <span class="tech-desc">RSI、MACD、MA等</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">预测引擎</span>
              <span class="tech-desc">智能预测系统</span>
            </div>
          </div>
        </div>
        
        <div class="tech-category deploy-tech">
          <div class="category-header">
            <div class="category-icon">🚀</div>
            <h4 class="category-title">部署运维</h4>
          </div>
          <div class="tech-list">
            <div class="tech-item">
              <span class="tech-name">Docker</span>
              <span class="tech-desc">容器化部署</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">Kubernetes</span>
              <span class="tech-desc">容器编排管理</span>
            </div>
            <div class="tech-item">
              <span class="tech-name">Nginx</span>
              <span class="tech-desc">负载均衡代理</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// About 页面逻辑
console.log('About page - Modern minimalist style')
</script>

<style lang="scss" scoped>
@use '../assets/styles/modern.scss' as *;

.modern-about {
  max-width: 1200px;
  margin: 0 auto;
}

// 页面标题
.page-header {
  text-align: center;
  margin-bottom: var(--claude-space-2xl);
  
  .page-title {
    font-size: 2.5rem;
    font-weight: 700;
    @include claude-text-gradient(var(--claude-primary), var(--claude-accent));
    margin-bottom: var(--claude-space);
    letter-spacing: -0.025em;
    
    .about-logo {
      font-weight: 600;
      margin-left: 0.5rem;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      
      /* 迷你狗狗logo */
      .mini-dog {
        position: relative;
        width: 24px;
        height: 24px;
        
        .mini-face {
          width: 20px;
          height: 20px;
          background: linear-gradient(135deg, var(--claude-primary-light), var(--claude-primary));
          border-radius: 50% 50% 60% 60%;
          position: absolute;
          top: 2px;
          left: 2px;
          
          &::before {
            content: '• •';
            position: absolute;
            top: 6px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 6px;
            color: var(--claude-text-primary);
            letter-spacing: 4px;
          }
          
          &::after {
            content: '';
            position: absolute;
            bottom: 4px;
            left: 50%;
            transform: translateX(-50%);
            width: 4px;
            height: 3px;
            background: var(--claude-accent);
            border-radius: 50%;
          }
        }
        
        .mini-ear-left {
          position: absolute;
          top: 0;
          left: 4px;
          width: 6px;
          height: 10px;
          background: var(--claude-primary);
          border-radius: 6px 6px 0 0;
          transform: rotate(-15deg);
        }
        
        .mini-ear-right {
          position: absolute;
          top: 0;
          right: 4px;
          width: 6px;
          height: 10px;
          background: var(--claude-primary);
          border-radius: 6px 6px 0 0;
          transform: rotate(15deg);
        }
      }
      
      .logo-text {
        font-weight: 700;
        background: linear-gradient(135deg, var(--claude-primary), var(--claude-accent));
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }
    
    @media (max-width: 768px) {
      font-size: 2rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1.75rem;
    }
  }
  
  .page-subtitle {
    font-size: 1.25rem;
    color: var(--claude-text-secondary);
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
    
    @media (max-width: 768px) {
      font-size: 1.125rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1rem;
      padding: 0 var(--claude-space-sm);
    }
  }
}

// 功能特性区域
.features-section {
  margin-bottom: var(--claude-space-2xl);
  
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: var(--claude-space-xl);
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
      gap: var(--claude-space-lg);
    }
  }
}

.feature-card {
  @include claude-card;
  padding: var(--claude-space-xl);
  text-align: center;
  transition: var(--claude-transition);
  animation: claude-fade-in 0.6s ease-out;
  
  &:hover {
    transform: translateY(-8px);
    box-shadow: var(--claude-shadow-xl);
  }
  
  &.ai-feature:hover {
    border-color: var(--claude-primary);
  }
  
  &.data-feature:hover {
    border-color: var(--claude-success);
  }
  
  &.service-feature:hover {
    border-color: var(--claude-accent);
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space-lg);
  }
}

.feature-header {
  margin-bottom: var(--claude-space-lg);
  
  .feature-icon {
    margin-bottom: var(--claude-space);
    
    .ai-brain {
      position: relative;
      width: 60px;
      height: 60px;
      margin: 0 auto;
      
      .brain-core {
        width: 100%;
        height: 100%;
        background: linear-gradient(135deg, var(--claude-primary), var(--claude-primary-light));
        border-radius: 50%;
        position: relative;
        
        &::after {
          content: '';
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 30px;
          height: 30px;
          background: var(--claude-bg-card);
          border-radius: 50%;
        }
      }
      
      .brain-waves {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        
        .wave {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          border: 2px solid var(--claude-primary);
          border-radius: 50%;
          animation: brain-pulse 2s infinite;
          
          &:nth-child(1) {
            width: 70px;
            height: 70px;
            animation-delay: 0s;
          }
          &:nth-child(2) {
            width: 90px;
            height: 90px;
            animation-delay: 0.5s;
          }
          &:nth-child(3) {
            width: 110px;
            height: 110px;
            animation-delay: 1s;
          }
        }
      }
    }
    
    .data-chart {
      position: relative;
      width: 60px;
      height: 40px;
      margin: 0 auto;
      
      .chart-line {
        position: absolute;
        top: 20px;
        left: 0;
        right: 0;
        height: 2px;
        background: var(--claude-success);
        
        &::before {
          content: '';
          position: absolute;
          left: 20px;
          top: -8px;
          width: 2px;
          height: 16px;
          background: var(--claude-success);
        }
        
        &::after {
          content: '';
          position: absolute;
          right: 15px;
          top: -15px;
          width: 2px;
          height: 30px;
          background: var(--claude-success);
        }
      }
      
      .data-points {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        
        .point {
          position: absolute;
          width: 6px;
          height: 6px;
          background: var(--claude-accent);
          border-radius: 50%;
          animation: data-pulse 2s infinite;
          
          &:nth-child(1) {
            top: 17px;
            left: 5px;
          }
          &:nth-child(2) {
            top: 12px;
            left: 20px;
            animation-delay: 0.3s;
          }
          &:nth-child(3) {
            top: 5px;
            right: 15px;
            animation-delay: 0.6s;
          }
          &:nth-child(4) {
            top: 17px;
            right: 5px;
            animation-delay: 0.9s;
          }
        }
      }
    }
    
    .service-monitor {
      width: 60px;
      height: 40px;
      margin: 0 auto;
      
      .monitor-screen {
        @include claude-card;
        width: 100%;
        height: 100%;
        padding: 8px;
        background: linear-gradient(135deg, var(--claude-accent), var(--claude-secondary));
        
        .screen-content {
          background: var(--claude-bg-card);
          width: 100%;
          height: 100%;
          border-radius: 2px;
          padding: 4px;
          display: flex;
          flex-direction: column;
          gap: 2px;
          
          .screen-line {
            height: 2px;
            background: var(--claude-success);
            border-radius: 1px;
            animation: screen-flicker 3s infinite;
            
            &:nth-child(1) {
              width: 80%;
              animation-delay: 0s;
            }
            &:nth-child(2) {
              width: 60%;
              animation-delay: 1s;
            }
            &:nth-child(3) {
              width: 90%;
              animation-delay: 2s;
            }
          }
        }
      }
    }
  }
  
  .feature-title {
    font-size: 1.375rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    margin: 0;
    
    @media (max-width: 480px) {
      font-size: 1.25rem;
    }
  }
}

.feature-content {
  .feature-description {
    font-size: 1rem;
    color: var(--claude-text-secondary);
    line-height: 1.6;
    margin-bottom: var(--claude-space-lg);
    
    @media (max-width: 480px) {
      font-size: 0.9rem;
    }
  }
  
  .tech-highlights {
    display: flex;
    flex-wrap: wrap;
    gap: var(--claude-space-sm);
    justify-content: center;
    
    .highlight-tag {
      background: var(--claude-bg-tertiary);
      color: var(--claude-text-secondary);
      padding: 0.375rem 0.75rem;
      border-radius: var(--claude-radius-lg);
      font-size: 0.75rem;
      font-weight: 500;
      border: 1px solid var(--claude-border-light);
      
      @media (max-width: 480px) {
        font-size: 0.7rem;
        padding: 0.3rem 0.6rem;
      }
    }
  }
}

// 风险提示
.risk-warning {
  @include claude-card;
  padding: var(--claude-space-xl);
  margin-bottom: var(--claude-space-2xl);
  border-left: 4px solid var(--claude-danger);
  background: rgba(239, 68, 68, 0.02);
  
  .warning-header {
    display: flex;
    align-items: center;
    gap: var(--claude-space);
    margin-bottom: var(--claude-space-lg);
    
    .warning-icon {
      .warning-symbol {
        font-size: 2rem;
        filter: drop-shadow(0 2px 4px rgba(239, 68, 68, 0.2));
      }
    }
    
    .warning-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: var(--claude-danger);
      margin: 0;
      
      @media (max-width: 480px) {
        font-size: 1.25rem;
      }
    }
  }
  
  .warning-content {
    .warning-text {
      font-size: 1rem;
      color: var(--claude-text-primary);
      line-height: 1.6;
      margin-bottom: var(--claude-space-lg);
      
      strong {
        color: var(--claude-danger);
        font-weight: 600;
      }
      
      @media (max-width: 480px) {
        font-size: 0.9rem;
      }
    }
    
    .warning-points {
      display: flex;
      flex-direction: column;
      gap: var(--claude-space-sm);
      
      .warning-point {
        display: flex;
        align-items: flex-start;
        gap: var(--claude-space-sm);
        
        .point-icon {
          color: var(--claude-danger);
          font-weight: bold;
          margin-top: 0.125rem;
        }
        
        .point-text {
          font-size: 0.875rem;
          color: var(--claude-text-secondary);
          line-height: 1.5;
          
          @media (max-width: 480px) {
            font-size: 0.8rem;
          }
        }
      }
    }
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space-lg);
  }
}

// 技术栈区域
.tech-stack-section {
  .section-title {
    font-size: 1.875rem;
    font-weight: 600;
    color: var(--claude-text-primary);
    text-align: center;
    margin-bottom: var(--claude-space-xl);
    
    @media (max-width: 480px) {
      font-size: 1.5rem;
    }
  }
  
  .tech-categories {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: var(--claude-space-lg);
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }
}

.tech-category {
  @include claude-card;
  padding: var(--claude-space-lg);
  transition: var(--claude-transition);
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--claude-shadow-lg);
  }
  
  &.frontend-tech:hover {
    border-color: var(--claude-success);
  }
  
  &.backend-tech:hover {
    border-color: var(--claude-accent);
  }
  
  &.ai-tech:hover {
    border-color: var(--claude-primary);
  }
  
  &.deploy-tech:hover {
    border-color: var(--claude-warning);
  }
  
  .category-header {
    display: flex;
    align-items: center;
    gap: var(--claude-space);
    margin-bottom: var(--claude-space-lg);
    
    .category-icon {
      font-size: 1.5rem;
    }
    
    .category-title {
      font-size: 1.125rem;
      font-weight: 600;
      color: var(--claude-text-primary);
      margin: 0;
      
      @media (max-width: 480px) {
        font-size: 1rem;
      }
    }
  }
  
  .tech-list {
    display: flex;
    flex-direction: column;
    gap: var(--claude-space);
    
    .tech-item {
      display: flex;
      flex-direction: column;
      gap: var(--claude-space-xs);
      
      .tech-name {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--claude-text-primary);
        
        @media (max-width: 480px) {
          font-size: 0.8rem;
        }
      }
      
      .tech-desc {
        font-size: 0.75rem;
        color: var(--claude-text-tertiary);
        line-height: 1.4;
        
        @media (max-width: 480px) {
          font-size: 0.7rem;
        }
      }
    }
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space);
  }
}

// 动画
@keyframes brain-pulse {
  0%, 100% {
    opacity: 0.3;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: 0.1;
    transform: translate(-50%, -50%) scale(1.1);
  }
}

@keyframes data-pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.3);
  }
}

@keyframes screen-flicker {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>