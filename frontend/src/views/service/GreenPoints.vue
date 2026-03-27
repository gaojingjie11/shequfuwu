<template>
  <div class="green-points-page">
    <Navbar />

    <div class="container">
      <div class="hero card">
        <div>
          <h1 class="page-title">绿色积分中心</h1>
          <p class="hero-desc">
            上传垃圾分类图片，AI 自动识别并发放积分。积分按 {{ GREEN_POINTS_PER_YUAN }} 积分 = 1 元，
            支付商城订单和物业费时会优先抵扣。
          </p>
        </div>
        <div class="hero-stats">
          <div class="stat">
            <span class="stat-label">当前积分余额</span>
            <span class="stat-value">{{ userStore.userInfo.green_points || 0 }}</span>
          </div>
          <div class="stat">
            <span class="stat-label">账户余额</span>
            <span class="stat-value">¥{{ formatAmount(userStore.userInfo.balance || 0) }}</span>
          </div>
        </div>
      </div>

      <div class="content-grid">
        <el-card class="upload-card">
          <template #header>
            <div class="card-header">AI 垃圾分类识别</div>
          </template>

          <el-upload
            drag
            :auto-upload="false"
            :show-file-list="false"
            accept="image/*"
            :on-change="handleFileChange"
          >
            <el-icon class="upload-icon"><UploadFilled /></el-icon>
            <div class="el-upload__text">拖拽图片到此处，或点击选择文件</div>
            <template #tip>
              <div class="el-upload__tip">建议上传清晰的垃圾分类照片，支持 jpg/png/webp</div>
            </template>
          </el-upload>

          <div v-if="previewUrl" class="preview-box">
            <img :src="previewUrl" alt="preview" class="preview-image" />
          </div>

          <div class="actions">
            <el-button type="success" :loading="uploading" :disabled="!selectedFile" @click="submitGarbageImage">
              {{ uploading ? '识别中...' : '开始识别并领取积分' }}
            </el-button>
          </div>

          <el-result
            v-if="recognitionResult"
            icon="success"
            :title="`本次奖励 ${recognitionResult.points} 积分`"
            :sub-title="recognitionResult.reason"
          >
            <template #extra>
              <el-tag type="success">当前积分余额：{{ recognitionResult.green_points }}</el-tag>
            </template>
          </el-result>
        </el-card>

        <el-card class="leaderboard-card">
          <template #header>
            <div class="card-header">累计绿色积分排行榜</div>
          </template>

          <el-table :data="leaderboard" stripe>
            <el-table-column label="排名" width="70">
              <template #default="{ row }">
                <el-tag :type="row.rank <= 3 ? 'success' : 'info'">#{{ row.rank }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="用户">
              <template #default="{ row }">
                <div class="user-cell">
                  <el-avatar :src="row.avatar" />
                  <span>{{ row.nickname || row.username }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="points" label="积分" width="100" />
          </el-table>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import Navbar from '@/components/layout/Navbar.vue'
import { uploadGarbageImage, getGreenPointsLeaderboard } from '@/api/greenPoints'
import { useUserStore } from '@/stores/user'
import { GREEN_POINTS_PER_YUAN } from '@/utils/payment'

const userStore = useUserStore()
const selectedFile = ref(null)
const previewUrl = ref('')
const uploading = ref(false)
const recognitionResult = ref(null)
const leaderboard = ref([])

function formatAmount(value) {
  return Number(value || 0).toFixed(2)
}

function handleFileChange(file) {
  selectedFile.value = file.raw
  recognitionResult.value = null
  if (file.raw) {
    previewUrl.value = URL.createObjectURL(file.raw)
  }
}

async function submitGarbageImage() {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择图片')
    return
  }

  const formData = new FormData()
  formData.append('file', selectedFile.value)

  uploading.value = true
  try {
    const res = await uploadGarbageImage(formData)
    recognitionResult.value = res
    ElMessage.success(`识别成功，奖励 ${res.points} 积分`)
    await Promise.all([fetchLeaderboard(), userStore.fetchUserInfo()])
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || error.message || '识别失败')
  } finally {
    uploading.value = false
  }
}

async function fetchLeaderboard() {
  const res = await getGreenPointsLeaderboard({ limit: 10 })
  leaderboard.value = res.list || []
}

onMounted(async () => {
  await Promise.all([fetchLeaderboard(), userStore.fetchUserInfo()])
})
</script>

<style scoped>
.green-points-page {
  min-height: 100vh;
  padding-bottom: 40px;
}

.hero {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.hero-desc {
  max-width: 720px;
  color: var(--text-secondary);
}

.hero-stats {
  display: flex;
  gap: 16px;
}

.stat {
  min-width: 140px;
  padding: 16px;
  border-radius: 12px;
  background: rgba(33, 150, 83, 0.08);
}

.stat-label {
  display: block;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #2e7d32;
}

.content-grid {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 24px;
}

.card-header {
  font-weight: 600;
}

.preview-box {
  margin: 16px 0;
  text-align: center;
}

.preview-image {
  max-width: 100%;
  max-height: 280px;
  border-radius: 12px;
}

.actions {
  margin: 16px 0 8px;
  text-align: center;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

@media (max-width: 960px) {
  .hero,
  .content-grid {
    grid-template-columns: 1fr;
    display: grid;
  }

  .hero-stats {
    flex-wrap: wrap;
  }
}
</style>
