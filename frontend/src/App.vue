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
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #ffffff;
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

// 头部
.header {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  
  .container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 60px;
    
    @media (max-width: 768px) {
      height: 56px;
    }
  }
  
  .title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #2c3e50;
    
    @media (max-width: 768px) {
      font-size: 1.25rem;
    }
  }
  
  .nav {
    display: flex;
    gap: 2rem;
    
    @media (max-width: 768px) {
      gap: 1rem;
    }
  }
  
  .nav-item {
    text-decoration: none;
    color: #666;
    font-weight: 500;
    padding: 0.5rem 0;
    border-bottom: 2px solid transparent;
    transition: all 0.2s ease;
    
    &:hover {
      color: #007bff;
    }
    
    &.router-link-active {
      color: #007bff;
      border-bottom-color: #007bff;
    }
    
    @media (max-width: 768px) {
      font-size: 0.9rem;
      padding: 0.25rem 0;
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

// 页脚
.footer {
  background: white;
  border-top: 1px solid #e9ecef;
  padding: 1rem 0;
  text-align: center;
  color: #6c757d;
  font-size: 0.875rem;
  
  @media (max-width: 768px) {
    padding: 0.75rem 0;
    font-size: 0.8rem;
  }
}
</style>
