<template>
  <div id="app">
    <el-container class="app-container">
      <!-- 头部导航 -->
      <el-header class="app-header">
        <div class="header-content">
          <div class="logo-section">
            <img src="/logo.svg" alt="智投预测" class="logo" />
            <h1 class="title">智投预测</h1>
            <span class="subtitle">AI股市指数预测平台</span>
          </div>
          
          <div class="nav-section">
            <el-menu
              :default-active="$route.path"
              class="nav-menu"
              mode="horizontal"
              router
            >
              <el-menu-item index="/">
                <el-icon><TrendCharts /></el-icon>
                <span>预测大盘</span>
              </el-menu-item>
              <el-menu-item index="/history">
                <el-icon><DataAnalysis /></el-icon>
                <span>历史数据</span>
              </el-menu-item>
              <el-menu-item index="/about">
                <el-icon><InfoFilled /></el-icon>
                <span>关于我们</span>
              </el-menu-item>
            </el-menu>
          </div>
        </div>
      </el-header>

      <!-- 主体内容 -->
      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>

      <!-- 页脚 -->
      <el-footer class="app-footer">
        <div class="footer-content">
          <p>&copy; 2024 智投预测. All rights reserved.</p>
          <p class="disclaimer">
            <el-icon><WarningFilled /></el-icon>
            风险提示: 本平台基于AI技术提供指数预测，结果仅供参考，不构成投资建议。股市有风险，投资需谨慎。
          </p>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import { TrendCharts, DataAnalysis, InfoFilled, WarningFilled } from '@element-plus/icons-vue'

export default {
  name: 'App',
  components: {
    TrendCharts,
    DataAnalysis,
    InfoFilled,
    WarningFilled
  }
}
</script>

<style lang="scss">
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.app-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 0;
  
  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    padding: 0 20px;
  }
  
  .logo-section {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .logo {
      width: 40px;
      height: 40px;
    }
    
    .title {
      margin: 0;
      font-size: 24px;
      font-weight: bold;
      color: #2c3e50;
    }
    
    .subtitle {
      font-size: 14px;
      color: #7f8c8d;
    }
  }
  
  .nav-menu {
    border: none;
    background: transparent;
    
    .el-menu-item {
      color: #2c3e50;
      
      &:hover, &.is-active {
        color: #409eff;
        background: rgba(64, 158, 255, 0.1);
      }
    }
  }
}

.app-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: calc(100vh - 120px - 80px);
}

.app-footer {
  background: rgba(0, 0, 0, 0.8);
  color: white;
  text-align: center;
  padding: 20px;
  
  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    
    p {
      margin: 5px 0;
    }
    
    .disclaimer {
      font-size: 12px;
      color: #ffd93d;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 5px;
    }
  }
}

// 路由过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
