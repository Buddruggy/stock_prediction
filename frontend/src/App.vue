<template>
  <div id="app">
    <!-- 头部 -->
    <header class="header">
      <div class="container">
        <h1 class="title">智投预测</h1>
        <nav class="nav">
          <router-link to="/" class="nav-item">预测</router-link>
          <router-link to="/history" class="nav-item">历史</router-link>
          <router-link to="/about" class="nav-item">关于</router-link>
        </nav>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <p>&copy; 2024 智投预测 | 投资需谨慎</p>
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
// 全局重置
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.5;
  color: #1a1a1a;
  background-color: #FBF9F6;
  font-size: 16px;
  font-weight: 400;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

// 容器
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
  
  @media (max-width: 768px) {
    padding: 0 0.75rem;
  }
}

// 头部 - 参考Claude官网风格
.header {
  background: #FBF9F6;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
  
  .container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    
    @media (max-width: 768px) {
      height: 56px;
    }
  }
  
  .title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a1a1a;
    letter-spacing: -0.025em;
    
    @media (max-width: 768px) {
      font-size: 1.25rem;
    }
  }
  
  .nav {
    display: flex;
    gap: 2rem;
    
    @media (max-width: 768px) {
      gap: 1.5rem;
    }
  }
  
  .nav-item {
    text-decoration: none;
    color: #6b7280;
    font-weight: 500;
    font-size: 0.95rem;
    padding: 0.5rem 0;
    transition: color 0.2s ease;
    
    &:hover {
      color: #2563eb;
    }
    
    &.router-link-active {
      color: #2563eb;
      font-weight: 600;
    }
    
    @media (max-width: 768px) {
      font-size: 0.9rem;
    }
  }
}

// 主体内容
.main {
  flex: 1;
  padding: 2rem 0;
  
  @media (max-width: 768px) {
    padding: 1rem 0;
  }
}

// 页脚 - 参考Claude官网风格
.footer {
  background: #FBF9F6;
  border-top: 1px solid #e5e7eb;
  padding: 1.5rem 0;
  text-align: center;
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 400;
  
  @media (max-width: 768px) {
    padding: 1rem 0;
    font-size: 0.8rem;
  }
}
</style>
