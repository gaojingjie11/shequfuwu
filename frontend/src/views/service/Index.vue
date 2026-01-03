<template>
  <div class="service-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">社区服务</h1>
      
      <div class="service-menu">
        <div class="service-card card" @click="$router.push('/service/notice')">
          <div class="service-icon">📢</div>
          <div class="service-name">公告通知</div>
          <div class="service-desc">查看社区最新公告</div>
        </div>
        
        <div class="service-card card" @click="$router.push('/service/repair')">
          <div class="service-icon">🔧</div>
          <div class="service-name">报修投诉</div>
          <div class="service-desc">提交报修和投诉</div>
        </div>
        
        <div class="service-card card" @click="$router.push('/service/visitor')">
          <div class="service-icon">👥</div>
          <div class="service-name">访客登记</div>
          <div class="service-desc">登记访客信息</div>
        </div>
        
        <div class="service-card card" @click="$router.push('/service/parking')">
          <div class="service-icon">🅿️</div>
          <div class="service-name">车位管理</div>
          <div class="service-desc">查看和绑定车位</div>
        </div>
        
        <div class="service-card card" @click="$router.push('/service/property')">
          <div class="service-icon">💰</div>
          <div class="service-name">物业费</div>
          <div class="service-desc">查看和缴纳物业费</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Navbar from '@/components/layout/Navbar.vue'
</script>

<style scoped>
.service-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.service-menu {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  margin-top: var(--spacing-xl);
}

.service-card {
  text-align: center;
  padding: var(--spacing-xl);
  cursor: pointer;
  transition: all var(--transition-base);
}

.service-card:hover {
  transform: translateY(-4px);
}

.service-icon {
  font-size: 64px;
  margin-bottom: var(--spacing-md);
}

.service-name {
  font-size: var(--font-size-xl);
  font-weight: 600;
  margin-bottom: var(--spacing-sm);
}

.service-desc {
  color: var(--text-secondary);
}
</style>
