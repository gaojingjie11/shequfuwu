<template>
  <div class="profile-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">个人中心</h1>
      
      <el-card class="profile-card">
        <div class="profile-header">
          <div class="avatar-container" @click="openEditProfile">
             <img :src="userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" class="avatar-img">
             <div class="avatar-overlay">修改头像</div>
          </div>
          <div class="user-info">
            <h2>{{ userInfo.real_name || userInfo.username }}</h2>
            <p>{{ userInfo.mobile }}</p>
            <el-tag size="small" effect="plain" style="margin-top: 5px">
                 {{ {'admin': '系统管理员', 'store': '商户', 'property': '物业'}[userInfo.role] || '居民' }}
            </el-tag>
          </div>
          <div class="header-actions">
              <el-button type="primary" plain @click="openEditProfile">编辑资料 / 修改密码</el-button>
          </div>
        </div>
        
        <div class="profile-stats">
          <div class="stat-item">
            <div class="stat-label">余额</div>
            <div class="stat-value">¥{{ userInfo.balance || 0 }}</div>
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
         <el-card shadow="hover" class="menu-item" @click="$router.push('/order')">
             <div class="menu-icon">📦</div>
             <span>我的订单</span>
         </el-card>
         <el-card shadow="hover" class="menu-item" @click="$router.push('/user/favorites')">
             <div class="menu-icon">❤️</div>
             <span>我的收藏</span>
         </el-card>
         <el-card shadow="hover" class="menu-item" @click="$router.push('/user/transactions')">
             <div class="menu-icon">💰</div>
             <span>我的账单</span>
         </el-card>
          <el-card shadow="hover" class="menu-item" @click="$router.push('/service')">
             <div class="menu-icon">🏘️</div>
             <span>社区服务</span>
         </el-card>
          <el-card shadow="hover" class="menu-item" @click="$router.push('/admin')" v-if="['admin', 'store', 'property'].includes(userInfo.role)">
             <div class="menu-icon">⚙️</div>
             <span>管理后台</span>
         </el-card>
          <el-card shadow="hover" class="menu-item" @click="handleLogout">
             <div class="menu-icon">🚪</div>
             <span>退出登录</span>
         </el-card>
      </div>
    </div>

    <!-- Consolidated Edit Profile Dialog -->
    <el-dialog v-model="showEditDialog" title="个人设置" width="500px">
        <el-tabs v-model="activeTab">
            <el-tab-pane label="基本信息" name="info">
                <el-form :model="editForm" label-width="80px">
                    <el-form-item label="头像">
                         <div class="avatar-uploader" @click="triggerFileUpload">
                             <img v-if="editForm.avatar" :src="editForm.avatar" class="upload-avatar">
                             <i v-else class="el-icon-plus avatar-uploader-icon">+</i>
                             <input type="file" ref="fileInput" @change="handleUpload" accept="image/*" style="display:none">
                         </div>
                    </el-form-item>
                    <el-form-item label="姓名">
                        <el-input v-model="editForm.real_name" />
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
                        <el-button type="primary" @click="submitInfo" :loading="loading">保存信息</el-button>
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
                        <el-button type="danger" @click="submitPwd" :loading="loading">确认修改</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
        </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { useUserStore } from '@/stores/user'
import { updateUserInfo, changePassword } from '@/api/user'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)

const showEditDialog = ref(false)
const activeTab = ref('info')
const loading = ref(false)
const fileInput = ref(null)

const editForm = ref({
    real_name: '',
    gender: 1,
    age: 0,
    email: '',
    avatar: ''
})

const pwdForm = ref({
    old_password: '',
    new_password: ''
})

const openEditProfile = () => {
    editForm.value = { ...userInfo.value }
    pwdForm.value = { old_password: '', new_password: '' }
    activeTab.value = 'info'
    showEditDialog.value = true
}

const triggerFileUpload = () => {
    fileInput.value.click()
}

const handleUpload = async (event) => {
  const file = event.target.files[0]
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
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

const submitInfo = async () => {
    loading.value = true
    try {
        await updateUserInfo(editForm.value)
        ElMessage.success('保存成功')
        await userStore.fetchUserInfo()
        showEditDialog.value = false
    } catch (e) {
        ElMessage.error(e.response?.data?.msg || '保存失败')
    } finally {
        loading.value = false
    }
}

const submitPwd = async () => {
    if (!pwdForm.value.old_password || !pwdForm.value.new_password) {
        ElMessage.warning('请输入密码')
        return
    }
    loading.value = true
    try {
        await changePassword(pwdForm.value)
        ElMessage.success('修改成功，请重新登录')
        userStore.logout()
        router.push('/login')
    } catch (e) {
        ElMessage.error(e.response?.data?.msg || '修改失败')
    } finally {
        loading.value = false
    }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
  }).then(() => {
    userStore.logout()
    router.push('/home')
  }).catch(() => {})
}

onMounted(async () => {
  await userStore.fetchUserInfo()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.profile-card {
    margin-bottom: 20px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding-bottom: var(--spacing-lg);
  position: relative;
}

.header-actions {
    margin-left: auto;
}

.avatar-container {
    position: relative;
    width: 80px;
    height: 80px;
    cursor: pointer;
    border-radius: 50%;
    overflow: hidden;
}

.avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.avatar-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    background: rgba(0,0,0,0.5);
    color: white;
    font-size: 10px;
    text-align: center;
    padding: 2px 0;
    opacity: 0;
    transition: opacity 0.3s;
}

.avatar-container:hover .avatar-overlay {
    opacity: 1;
}

.profile-stats {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #409EFF;
}

.menu-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 15px;
}

.menu-item {
    cursor: pointer;
    text-align: center;
    transition: transform 0.2s;
}
.menu-item:hover {
    transform: translateY(-5px);
}

.menu-icon {
    font-size: 28px;
    margin-bottom: 10px;
}

/* Avatar Uploader in Dialog */
.avatar-uploader {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    width: 100px;
    height: 100px;
    display: flex;
    justify-content: center;
    align-items: center;
}
.avatar-uploader:hover {
    border-color: #409EFF;
}
.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
}
.upload-avatar {
    width: 100px;
    height: 100px;
    display: block;
    object-fit: cover;
}
</style>
