<template>
  <div class="navbar">
    <div class="container">
      <div class="navbar-content">
        <div class="navbar-logo" @click="$router.push('/home')">
          🏘️ 智慧社区
        </div>
        
        <div class="navbar-menu">
          <router-link to="/home" class="nav-link">首页</router-link>
          <router-link to="/mall" class="nav-link">商城</router-link>
          <router-link to="/service" class="nav-link">社区服务</router-link>
          <router-link to="/chat" class="nav-link" v-if="userStore.isLoggedIn">AI助手</router-link>
          <router-link to="/order" class="nav-link" v-if="userStore.isLoggedIn">我的订单</router-link>
        </div>
        
        <div class="navbar-right">
          <router-link to="/cart" class="cart-icon" v-if="userStore.isLoggedIn">
            🛒
            <span class="cart-badge" v-if="cartStore.cartCount > 0">
              {{ cartStore.cartCount }}
            </span>
          </router-link>
          
          <div v-if="userStore.isLoggedIn" class="user-menu">
            <div class="user-avatar" @click="showUserMenu = !showUserMenu" :title="userStore.userInfo.real_name || userStore.userInfo.username">
               <img :src="userStore.userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" style="width:100%; height:100%; border-radius:50%; object-fit:cover;">
            </div>
            <div class="user-dropdown" v-if="showUserMenu">
              <router-link to="/profile" class="dropdown-item">个人中心</router-link>
              <router-link to="/admin" class="dropdown-item" v-if="['admin', 'store', 'property'].includes(userStore.userInfo.role)">
                管理后台
              </router-link>
              <div class="dropdown-item" @click="handleLogout">退出登录</div>
            </div>
          </div>
          
          <div v-else class="auth-buttons">
            <router-link to="/login" class="btn btn-sm btn-outline">登录</router-link>
            <router-link to="/register" class="btn btn-sm btn-primary">注册</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useCartStore } from '@/stores/cart'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()
const showUserMenu = ref(false)

const handleLogout = () => {
  userStore.logout()
  showUserMenu.value = false
  router.push('/home')
}
</script>

<style scoped>
.navbar {
  background: var(--bg-white);
  box-shadow: var(--shadow-sm);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.navbar-logo {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--primary-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.navbar-menu {
  display: flex;
  gap: var(--spacing-xl);
}

.nav-link {
  color: var(--text-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color var(--transition-fast);
  position: relative;
}

.nav-link:hover, .nav-link.router-link-active {
  color: var(--primary-color);
}

.nav-link.router-link-active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--primary-color);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.cart-icon {
  position: relative;
  font-size: 24px;
  cursor: pointer;
  text-decoration: none;
  transition: transform var(--transition-fast);
}

.cart-icon:hover {
  transform: scale(1.1);
}

.cart-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: var(--danger-color);
  color: white;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.user-menu {
  position: relative;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-weight: 500;
  color: var(--primary-dark);
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: white;
  border-radius: var(--border-radius-sm);
  box-shadow: var(--shadow-md);
  min-width: 150px;
  overflow: hidden;
}

.dropdown-item {
  display: block;
  padding: 12px 16px;
  color: var(--text-primary);
  cursor: pointer;
  text-decoration: none;
  transition: background var(--transition-fast);
}

.dropdown-item:hover {
  background: var(--bg-gray);
}

.auth-buttons {
  display: flex;
  gap: var(--spacing-sm);
}
</style>
