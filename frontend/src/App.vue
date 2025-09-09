<template>
  <div id="app">
    <!-- Claude 风格头部 -->
    <header class="claude-header">
      <div class="claude-container">
        <div class="header-content">
          <div class="brand-section">
            <h1 class="brand-title">
              <div class="dog-logo">
                <div class="dog-head">
                  <!-- 狗狗的耳朵 -->
                  <div class="ear ear-left"></div>
                  <div class="ear ear-right"></div>
                  <!-- 狗狗的脸 -->
                  <div class="face">
                    <!-- 眼睛 -->
                    <div class="eye eye-left"></div>
                    <div class="eye eye-right"></div>
                    <!-- 鼻子 -->
                    <div class="nose"></div>
                  </div>
                </div>
              </div>
            </h1>
          </div>
          
          <nav class="main-nav">
            <router-link to="/" class="nav-link">
              <span class="nav-text">预测</span>
            </router-link>
            <router-link to="/history" class="nav-link">
              <span class="nav-text">历史</span>
            </router-link>
            <router-link to="/about" class="nav-link">
              <span class="nav-text">关于</span>
            </router-link>
          </nav>
        </div>
      </div>
    </header>

    <!-- Claude 风格主体内容 -->
    <main class="claude-main">
      <div class="claude-container">
        <router-view />
      </div>
    </main>

    <!-- Claude 风格页脚 -->
    <footer class="claude-footer">
      <div class="claude-container">
        <div class="footer-content">
          <div class="footer-info">
            <span class="copyright">&copy; 2024 <span class="footer-logo">GoGoTou</span></span>
            <span class="divider">•</span>
            <span class="disclaimer">投资需谨慎</span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style lang="scss">
@use './src/assets/styles/modern.scss' as *;

// 全局重置
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--claude-font-sans);
  line-height: 1.6;
  color: var(--claude-text-primary);
  background: var(--claude-bg-primary);
  font-size: 16px;
  font-weight: 400;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

// Claude 风格头部
.claude-header {
  background: var(--claude-bg-card);
  border-bottom: 1px solid var(--claude-border);
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  
  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.25rem 0;
    min-height: 35px;
    
    @media (max-width: 768px) {
      flex-direction: row; /* 移动端保持水平布局 */
      gap: var(--claude-space-sm);
      padding: 0.2rem 0;
      min-height: 30px;
    }
    
    @media (max-width: 480px) {
      min-height: 28px;
    }
  }
  
  .brand-section {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .brand-title {
    margin: 0;
    
    .dog-logo {
      display: flex;
      align-items: center;
      
      .dog-head {
        position: relative;
        width: 32px; /* 增大狗狗图标 */
        height: 32px;
        
        @media (max-width: 768px) {
          width: 28px;
          height: 28px;
        }
        
        @media (max-width: 480px) {
          width: 24px;
          height: 24px;
        }
        
        /* 狗狗耳朵 */
        .ear {
          position: absolute;
          width: 10px; /* 相应增大耳朵 */
          height: 16px;
          background: var(--claude-primary);
          border-radius: 8px 8px 0 0;
          
          &.ear-left {
            top: 1px;
            left: 7px;
            transform: rotate(-20deg);
            animation: ear-wiggle 4s ease-in-out infinite;
          }
          
          &.ear-right {
            top: 1px;
            right: 7px;
            transform: rotate(20deg);
            animation: ear-wiggle 4s ease-in-out infinite;
            animation-delay: 0.5s;
          }
        }
        
        /* 狗狗的脸 */
        .face {
          position: absolute;
          top: 8px;
          left: 3px;
          width: 26px; /* 相应增大脸部 */
          height: 26px;
          
          @media (max-width: 768px) {
            top: 7px;
            left: 3px;
            width: 22px;
            height: 22px;
          }
          
          @media (max-width: 480px) {
            top: 6px;
            left: 2px;
            width: 20px;
            height: 20px;
          }
          background: linear-gradient(135deg, var(--claude-primary-light), var(--claude-primary));
          border-radius: 50% 50% 60% 60%;
          
          /* 眼睛 */
          .eye {
            position: absolute;
            width: 5px; /* 增大眼睛 */
            height: 5px;
            background: var(--claude-text-primary);
            border-radius: 50%;
            top: 10px;
            animation: blink 3s ease-in-out infinite;
            
            &.eye-left {
              left: 6px;
            }
            
            &.eye-right {
              right: 6px;
              animation-delay: 0.1s;
            }
          }
          
          /* 鼻子 */
          .nose {
            position: absolute;
            bottom: 6px;
            left: 50%;
            transform: translateX(-50%);
            width: 7px; /* 增大鼻子 */
            height: 5px;
            background: var(--claude-accent);
            border-radius: 50% 50% 50% 50% / 60% 60% 40% 40%;
            
            &::after {
              content: '';
              position: absolute;
              top: 5px;
              left: 50%;
              transform: translateX(-50%);
              width: 1px;
              height: 4px;
              background: var(--claude-accent);
              border-radius: 0 0 1px 1px;
            }
          }
        }
        
        /* 悬停效果 */
        &:hover {
          transform: scale(1.1) rotate(5deg);
          transition: transform 0.3s ease;
        }
      }
    }
    
    @media (max-width: 768px) {
      font-size: 1.75rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1.5rem;
    }
  }

  
  .main-nav {
    display: flex;
    gap: var(--claude-space-sm);
    
    @media (max-width: 768px) {
      gap: var(--claude-space-xs);
    }
  }
  
  .nav-link {
    @include claude-button-secondary;
    text-decoration: none;
    padding: 0.25rem 0.5rem; /* 缩小按钮内边距 */
    font-size: 0.8rem; /* 缩小字体 */
    font-weight: 500;
    transition: var(--claude-transition);
    
    &:hover {
      background: var(--claude-bg-hover);
      border-color: var(--claude-primary);
      color: var(--claude-primary);
      transform: translateY(-1px);
    }
    
    &.router-link-active {
      @include claude-button-primary;
      padding: 0.25rem 0.5rem; /* 覆盖默认的大内边距 */
      font-size: 0.8rem; /* 保持一致的字体大小 */
      
      @media (max-width: 768px) {
        padding: 0.2rem 0.4rem;
        font-size: 0.75rem;
      }
      
      @media (max-width: 480px) {
        padding: 0.15rem 0.35rem;
        font-size: 0.7rem;
      }
    }
    
    @media (max-width: 768px) {
      padding: 0.2rem 0.4rem;
      font-size: 0.75rem;
    }
    
    @media (max-width: 480px) {
      padding: 0.15rem 0.35rem;
      font-size: 0.7rem;
    }
  }
}

// Claude 风格主体内容
.claude-main {
  flex: 1;
  padding: var(--claude-space-lg) 0;
  animation: claude-fade-in 0.6s ease-out;
  
  @media (max-width: 768px) {
    padding: var(--claude-space-lg) 0;
  }
  
  @media (max-width: 480px) {
    padding: var(--claude-space) 0;
  }
}

// Claude 风格页脚
.claude-footer {
  background: var(--claude-bg-tertiary);
  border-top: 1px solid var(--claude-border);
  
  .footer-content {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--claude-space-lg) 0;
    
    @media (max-width: 768px) {
      padding: var(--claude-space) 0;
    }
  }
  
  .footer-info {
    display: flex;
    align-items: center;
    gap: var(--claude-space);
    font-size: 0.9rem;
    color: var(--claude-text-secondary);
    
    .divider {
      color: var(--claude-text-tertiary);
    }
    
    .disclaimer {
      color: var(--claude-warning);
      font-weight: 500;
    }
    
    .footer-logo {
      font-weight: 600;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      font-size: 0.95rem;
      
      /* 简化版狗狗logo */
      background: linear-gradient(90deg, var(--claude-primary) 0%, var(--claude-primary-light) 30%, var(--claude-accent) 70%, var(--claude-success) 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      
      &::before {
        content: '🐶';
        margin-right: 0.3em;
        font-size: 0.8em;
      }
    }
    
    @media (max-width: 480px) {
      font-size: 0.85rem;
      gap: var(--claude-space-sm);
    }
  }

}

// 容器样式
.claude-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--claude-space-lg);
  
  @media (max-width: 768px) {
    padding: 0 var(--claude-space);
  }
  
  @media (max-width: 480px) {
    padding: 0 var(--claude-space-sm);
  }
}

// 狗狗Logo动画
@keyframes ear-wiggle {
  0%, 100% {
    transform: rotate(-20deg);
  }
  50% {
    transform: rotate(-25deg);
  }
}

.ear-right {
  @keyframes ear-wiggle {
    0%, 100% {
      transform: rotate(20deg);
    }
    50% {
      transform: rotate(25deg);
    }
  }
}

@keyframes blink {
  0%, 90%, 100% {
    transform: scaleY(1);
  }
  95% {
    transform: scaleY(0.1);
  }
}

@keyframes tail-wag {
  0%, 100% {
    transform: rotate(-10deg);
  }
  50% {
    transform: rotate(10deg);
  }
}
</style>