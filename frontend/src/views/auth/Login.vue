<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card card">
        <h2 class="login-title">欢迎回来</h2>
        <p class="login-subtitle">登录智慧社区</p>
        
        <div class="tabs">
            <div class="tab-item" :class="{active: loginType === 'password'}" @click="loginType = 'password'">密码登录</div>
            <div class="tab-item" :class="{active: loginType === 'code'}" @click="loginType = 'code'">验证码登录</div>
        </div>

        <form v-if="loginType === 'password'" @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label>手机号</label>
            <input v-model="form.mobile" type="tel" class="input" placeholder="请输入手机号" required />
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="form.password" type="password" class="input" placeholder="请输入密码" required />
          </div>
          <div class="form-footer">
            <router-link to="/register" class="link">还没有账号？注册</router-link>
          </div>
          <button type="submit" class="btn btn-primary btn-lg" :disabled="loading">
            <span v-if="!loading">登录</span>
            <span v-else class="loading"></span>
          </button>
        </form>

        <form v-else @submit.prevent="handleCodeLogin" class="login-form">
          <div class="form-group">
            <label>手机号</label>
            <input v-model="codeForm.mobile" type="tel" class="input" placeholder="请输入手机号" required />
          </div>
          <div class="form-group">
            <label>验证码</label>
            <div class="code-input-group">
                <input v-model="codeForm.code" type="text" class="input" placeholder="6位验证码" required />
                <button type="button" class="btn btn-secondary btn-code" :disabled="timer > 0" @click="sendSms">
                    {{ timer > 0 ? `${timer}s` : '发送验证码' }}
                </button>
            </div>
          </div>
          <div class="form-footer">
            <router-link to="/register" class="link">还没有账号？注册</router-link>
          </div>
          <button type="submit" class="btn btn-primary btn-lg" :disabled="loading">
            <span v-if="!loading">登录</span>
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
import { sendCode, loginByCode } from '@/api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginType = ref('password') // password | code
const loading = ref(false)
const timer = ref(0)

const form = ref({ mobile: '', password: '' })
const codeForm = ref({ mobile: '', code: '' })

// 密码登录
const handleLogin = async () => {
  loading.value = true
  try {
    await userStore.login(form.value)
    ElMessage.success('登录成功！')
    router.push('/home')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 发送验证码
const sendSms = async () => {
    if (!codeForm.value.mobile) {
        ElMessage.warning('请输入手机号')
        return
    }
    try {
        await sendCode({ mobile: codeForm.value.mobile })
        ElMessage.success('验证码已发送')
        timer.value = 60
        const interval = setInterval(() => {
            timer.value--
            if (timer.value <= 0) clearInterval(interval)
        }, 1000)
    } catch (e) {
        ElMessage.error(e.response?.data?.message || '发送失败')
    }
}

// 验证码登录
const handleCodeLogin = async () => {
    loading.value = true
    try {
        await userStore.loginByCode(codeForm.value)
        ElMessage.success('登录成功')
        router.push('/home')
        router.push('/home')
    } catch (e) {
        // request.js interceptor returns new Error(msg), so we check e.message
        ElMessage.error(e.message || e.response?.data?.msg || '登录失败')
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

/* Tabs */
.tabs {
    display: flex;
    justify-content: center;
    border-bottom: 1px solid #eee;
    margin-bottom: 20px;
}
.tab-item {
    padding: 10px 20px;
    cursor: pointer;
    color: #666;
    border-bottom: 2px solid transparent;
}
.tab-item.active {
    color: var(--primary-color);
    border-bottom-color: var(--primary-color);
    font-weight: bold;
}

/* Code Input */
.code-input-group {
    display: flex;
    gap: 10px;
}
.btn-code {
    white-space: nowrap;
    width: 120px;
}
</style>
