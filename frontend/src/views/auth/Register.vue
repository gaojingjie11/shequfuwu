<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card card">
        <h2 class="login-title">创建账号</h2>
        <p class="login-subtitle">加入智慧社区</p>
        
        <form @submit.prevent="handleRegister" class="login-form">
          <div class="form-group">
            <label>手机号</label>
            <input 
              v-model="form.mobile" 
              type="tel" 
              class="input" 
              placeholder="请输入手机号"
              required
            />
          </div>
          
          <div class="form-group">
            <label>真实姓名</label>
            <input 
              v-model="form.real_name" 
              type="text" 
              class="input" 
              placeholder="请输入真实姓名"
              required
            />
          </div>
          
          <div class="form-group">
            <label>密码</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="input" 
              placeholder="请输入密码（至少6位）"
              required
              minlength="6"
            />
          </div>
          
          <div class="form-footer">
            <router-link to="/login" class="link">已有账号？登录</router-link>
          </div>
          
          <button type="submit" class="btn btn-primary btn-lg" :disabled="loading">
            <span v-if="!loading">注册</span>
            <span v-else class="loading"></span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  mobile: '',
  real_name: '',
  password: '',
  username: '',
  age: 0,
  gender: 0
})

const loading = ref(false)

const handleRegister = async () => {
  // 自动生成username
  if (!form.value.username) {
    form.value.username = 'user_' + form.value.mobile.slice(-6)
  }
  
  loading.value = true
  try {
    await userStore.register(form.value)
    ElMessage.success('注册成功！请登录')
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-light) 0%, var(--secondary-light) 100%);
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: var(--spacing-lg);
}

.login-card {
  animation: fadeIn 0.5s ease;
}

.login-title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  text-align: center;
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
}

.login-subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xl);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.form-group label {
  font-weight: 500;
  color: var(--text-primary);
}

.form-footer {
  display: flex;
  justify-content: flex-end;
}

.link {
  color: var(--primary-color);
  text-decoration: none;
  font-size: var(--font-size-sm);
}

.link:hover {
  text-decoration: underline;
}

.btn-lg {
  width: 100%;
}
</style>
