<template>
  <div class="repair-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">报修投诉</h1>
      
      <div class="repair-form card">
        <h3>提交新的报修/投诉</h3>
        <form @submit.prevent="submitRepair">
          <div class="form-group">
            <label>类型</label>
            <select v-model="form.type" class="input">
              <option :value="1">报修</option>
              <option :value="2">投诉</option>
            </select>
          </div>
          
          <div class="form-group">
            <label>分类</label>
            <input v-model="form.category" class="input" placeholder="如: 水电、门窗、噪音等" required />
          </div>
          
          <div class="form-group">
            <label>详细描述</label>
            <textarea v-model="form.content" class="input" rows="4" placeholder="请详细描述问题" required></textarea>
          </div>
          
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '提交中...' : '提交' }}
          </button>
        </form>
      </div>
      
      <h3 class="section-title">我的记录</h3>
      <div class="repair-list">
        <div class="repair-card card" v-for="item in repairs" :key="item.id">
          <div class="repair-header">
            <span class="tag" :class="item.type === 1 ? 'tag-info' : 'tag-warning'">
              {{ item.type === 1 ? '报修' : '投诉' }}
            </span>
            <span class="repair-category">{{ item.category }}</span>
            <span class="tag" :class="getStatusClass(item.status)">
              {{ getStatusText(item.status) }}
            </span>
          </div>
          <p class="repair-content">{{ item.content }}</p>
          <p class="repair-result" v-if="item.result">
            <strong>处理结果:</strong> {{ item.result }}
          </p>
          <div class="repair-time">{{ formatDate(item.created_at) }}</div>
        </div>
      </div>

      <div class="pagination-container mt-4" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { createRepair, getRepairList } from '@/api/service'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const form = ref({
  type: 1,
  category: '',
  content: ''
})

const loading = ref(false)
const repairs = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (status) => {
  const map = { 0: '待处理', 1: '处理中', 2: '已完成' }
  return map[status] || '未知'
}

const getStatusClass = (status) => {
  const map = { 0: 'tag-warning', 1: 'tag-info', 2: 'tag-success' }
  return map[status] || ''
}

const submitRepair = async () => {
  loading.value = true
  try {
    await createRepair(form.value)
    ElMessage.success('提交成功！')
    form.value = { type: 1, category: '', content: '' }
    await fetchRepairs()
  } catch (error) {
    ElMessage.error('提交失败')
  } finally {
    loading.value = false
  }
}

const fetchRepairs = async () => {
  try {
    const res = await getRepairList({
        page: currentPage.value,
        size: pageSize.value
    })
    repairs.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取记录失败:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchRepairs()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchRepairs()
}

onMounted(() => {
  fetchRepairs()
})
</script>

<style scoped>
.repair-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.repair-form {
  margin-bottom: var(--spacing-xl);
}

.repair-form h3 {
  margin-bottom: var(--spacing-lg);
}

.repair-form form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.section-title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  margin-bottom: var(--spacing-lg);
}

.repair-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.repair-card {
  padding: var(--spacing-lg);
}

.repair-header {
  display: flex;
  gap: var(--spacing-md);
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.repair-category {
  font-weight: 600;
}

.repair-content {
  margin-bottom: var(--spacing-md);
  line-height: 1.6;
}

.repair-result {
  padding: var(--spacing-md);
  background: var(--bg-gray);
  border-radius: var(--border-radius-sm);
  margin-bottom: var(--spacing-sm);
}

.repair-time {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

textarea.input {
  resize: vertical;
  font-family: inherit;
}
</style>
