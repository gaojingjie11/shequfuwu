<template>
  <div class="profile-page">
    <Navbar />

    <div class="container">
      <h1 class="page-title">个人中心</h1>

      <el-card class="profile-card">
        <div class="profile-header">
          <div class="avatar-container" @click="openEditProfile">
            <img :src="userInfo.avatar || defaultAvatar" class="avatar-img" />
            <div class="avatar-overlay">编辑资料</div>
          </div>

          <div class="user-info">
            <h2>{{ userInfo.real_name || userInfo.username }}</h2>
            <p>{{ userInfo.mobile }}</p>
            <el-tag size="small">{{ roleText }}</el-tag>
          </div>

          <div class="header-actions">
            <el-button type="primary" plain @click="openEditProfile">编辑资料</el-button>
          </div>
        </div>

        <div class="profile-stats">
          <div class="stat-item">
            <div class="stat-label">余额</div>
            <div class="stat-value">¥{{ formatAmount(userInfo.balance) }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">绿色积分</div>
            <div class="stat-value green">{{ userInfo.green_points || 0 }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">状态</div>
            <el-tag :type="userInfo.status === 1 ? 'success' : 'danger'">
              {{ userInfo.status === 1 ? '正常' : '冻结' }}
            </el-tag>
          </div>
        </div>
      </el-card>

      <div class="menu-grid">
        <el-card
          v-for="item in menuItems"
          :key="item.title"
          shadow="hover"
          class="menu-item"
          @click="item.action"
        >
          <div class="menu-icon-wrap">
            <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
          </div>
          <div class="menu-title">{{ item.title }}</div>
          <span>{{ item.desc }}</span>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="showEditDialog" title="个人设置" width="520px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="info">
          <el-form :model="editForm" label-width="80px">
            <el-form-item label="头像">
              <div class="avatar-uploader" @click="triggerFileUpload">
                <img v-if="editForm.avatar" :src="editForm.avatar" class="upload-avatar" />
                <span v-else>上传</span>
                <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="handleUpload" />
              </div>
            </el-form-item>
            <el-form-item label="姓名">
              <el-input v-model="editForm.real_name" />
            </el-form-item>
            <el-form-item label="昵称">
              <el-input v-model="editForm.username" />
            </el-form-item>
            <el-form-item label="性别">
              <el-radio-group v-model="editForm.gender">
                <el-radio :label="1">男</el-radio>
                <el-radio :label="2">女</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="年龄">
              <el-input-number v-model="editForm.age" :min="1" :max="120" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="editForm.email" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="submitInfo">保存信息</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="pwd">
          <el-form :model="pwdForm" label-width="80px">
            <el-form-item label="旧密码">
              <el-input v-model="pwdForm.old_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="pwdForm.new_password" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" :loading="loading" @click="submitPwd">确认修改</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ShoppingBag, Trophy, Star, Wallet, Service, Setting, SwitchButton } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'
import { updateUserInfo, changePassword } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo || {})
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const showEditDialog = ref(false)
const activeTab = ref('info')
const loading = ref(false)
const fileInput = ref(null)
const editForm = ref({})
const pwdForm = ref({ old_password: '', new_password: '' })

const roleText = computed(() => {
  const map = { admin: '系统管理员', store: '商户', property: '物业人员', user: '居民' }
  return map[userInfo.value.role] || '居民'
})

const menuItems = computed(() => {
  const items = [
    {
      title: '订单',
      desc: '我的订单',
      icon: ShoppingBag,
      action: () => router.push('/order')
    },
    {
      title: '积分',
      desc: '绿色积分中心',
      icon: Trophy,
      action: () => router.push('/service/green-points')
    },
    {
      title: '收藏',
      desc: '我的收藏',
      icon: Star,
      action: () => router.push('/user/favorites')
    },
    {
      title: '账单',
      desc: '交易流水',
      icon: Wallet,
      action: () => router.push('/user/transactions')
    },
    {
      title: '服务',
      desc: '社区服务',
      icon: Service,
      action: () => router.push('/service')
    }
  ]

  if (['admin', 'store', 'property'].includes(userInfo.value.role)) {
    items.push({
      title: '后台',
      desc: '管理后台',
      icon: Setting,
      action: () => router.push('/admin')
    })
  }

  items.push({
    title: '退出',
    desc: '退出登录',
    icon: SwitchButton,
    action: handleLogout
  })

  return items
})

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function openEditProfile() {
  editForm.value = { ...userInfo.value }
  pwdForm.value = { old_password: '', new_password: '' }
  activeTab.value = 'info'
  showEditDialog.value = true
}

function triggerFileUpload() {
  fileInput.value?.click()
}

async function handleUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await request({
      url: '/upload',
      method: 'post',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    editForm.value.avatar = res.url
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

async function submitInfo() {
  loading.value = true
  try {
    await updateUserInfo(editForm.value)
    await userStore.fetchUserInfo()
    ElMessage.success('保存成功')
    showEditDialog.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '保存失败')
  } finally {
    loading.value = false
  }
}

async function submitPwd() {
  if (!pwdForm.value.old_password || !pwdForm.value.new_password) {
    ElMessage.warning('请输入完整密码')
    return
  }

  loading.value = true
  try {
    await changePassword(pwdForm.value)
    ElMessage.success('密码已修改，请重新登录')
    await userStore.logout()
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '修改失败')
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  ElMessageBox.confirm('确认退出登录吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    await userStore.logout()
    router.push('/home')
  }).catch(() => {})
}

onMounted(() => {
  userStore.fetchUserInfo()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-container {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: auto 0 0;
  text-align: center;
  padding: 4px 0;
  font-size: 12px;
  color: #fff;
  background: rgba(0, 0, 0, 0.55);
}

.user-info {
  flex: 1;
}

.header-actions {
  margin-left: auto;
}

.profile-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 20px;
}

.stat-item {
  padding: 16px;
  border-radius: 12px;
  background: var(--bg-gray);
  text-align: center;
}

.stat-label {
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--primary-color);
}

.stat-value.green {
  color: #2e7d32;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 15px;
}

.menu-item {
  cursor: pointer;
  text-align: center;
  transition: all var(--transition-base);
}

.menu-item:hover {
  transform: translateY(-4px);
}

.menu-icon-wrap {
  width: 58px;
  height: 58px;
  border-radius: 16px;
  margin: 0 auto 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 180, 137, 0.12), rgba(0, 180, 137, 0.2));
}

.menu-icon {
  font-size: 30px;
  color: var(--primary-color);
}

.menu-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 6px;
}

.avatar-uploader {
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
}

.upload-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

@media (max-width: 768px) {
  .profile-header,
  .profile-stats {
    display: grid;
    grid-template-columns: 1fr;
  }
}
</style>
