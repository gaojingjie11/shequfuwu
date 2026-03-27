<template>
  <div class="visitor-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">访客登记</h1>
      
      <div class="visitor-form card">
        <h3>新增访客</h3>
        <form @submit.prevent="submitVisitor">
          <div class="form-group">
            <label>访客姓名</label>
            <input v-model="form.name" class="input" placeholder="请输入访客姓名" required />
          </div>
          
          <div class="form-group">
            <label>访客电话</label>
            <input v-model="form.mobile" type="tel" class="input" placeholder="请输入访客电话" required />
          </div>
          
          <div class="form-group">
            <label>来访原因</label>
            <input v-model="form.reason" class="input" placeholder="如: 探亲访友、快递配送等" required />
          </div>
          
          <div class="form-group">
            <label>预计来访时间</label>
            <input v-model="form.visit_time" type="datetime-local" class="input" required />
          </div>
          
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '提交中...' : '提交' }}
          </button>
        </form>
      </div>
      
      <h3 class="section-title">访客记录</h3>
      <div class="visitor-list">
        <div class="visitor-card card" v-for="visitor in visitors" :key="visitor.id">
          <div class="visitor-header">
            <div class="visitor-name">{{ visitor.name }}</div>
            <span class="tag" :class="getStatusClass(visitor.status)">
              {{ getStatusText(visitor.status) }}
            </span>
          </div>
          <div class="visitor-info">
            <div>电话: {{ visitor.mobile }}</div>
            <div>原因: {{ visitor.reason }}</div>
            <div>来访时间: {{ formatDate(visitor.visit_time) }}</div>
            <div v-if="visitor.audit_remark">备注: {{ visitor.audit_remark }}</div>
          </div>
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
import { createVisitor, getVisitorList } from '@/api/service'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const form = ref({
  name: '',
  mobile: '',
  reason: '',
  visit_time: ''
})

const loading = ref(false)
const visitors = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (status) => {
  const map = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
}

const getStatusClass = (status) => {
  const map = { 0: 'tag-warning', 1: 'tag-success', 2: 'tag-danger' }
  return map[status] || ''
}

const submitVisitor = async () => {
  loading.value = true
  try {
    // 格式化时间为后端需要的 YYYY-MM-DD HH:mm:ss
    const submitData = {
      ...form.value,
      visitor_name: form.value.name,
      visitor_phone: form.value.mobile,
      visit_time: dayjs(form.value.visit_time).format('YYYY-MM-DD HH:mm:ss')
    }
    
    await createVisitor(submitData)
    ElMessage.success('提交成功，等待审核！')
    form.value = { name: '', mobile: '', reason: '', visit_time: '' }
    await fetchVisitors()
  } catch (error) {
    ElMessage.error('提交失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

const fetchVisitors = async () => {
  try {
    const res = await getVisitorList({
        page: currentPage.value,
        size: pageSize.value
    })
    visitors.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取访客记录失败:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchVisitors()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchVisitors()
}

onMounted(() => {
  fetchVisitors()
})
</script>

<style scoped>
.visitor-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.visitor-form {
  margin-bottom: var(--spacing-xl);
}

.visitor-form h3 {
  margin-bottom: var(--spacing-lg);
}

.visitor-form form {
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

.visitor-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.visitor-card {
  padding: var(--spacing-lg);
}

.visitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.visitor-name {
  font-size: var(--font-size-lg);
  font-weight: 600;
}

.visitor-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  color: var(--text-secondary);
}
</style>
