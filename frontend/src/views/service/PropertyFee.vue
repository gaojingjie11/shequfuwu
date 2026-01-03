<template>
  <div class="property-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">物业费缴纳</h1>
      
      <div class="user-balance card">
        <span>当前余额：</span>
        <span class="balance-amount">¥{{ userStore.userInfo.balance || 0 }}</span>
      </div>
      
      <div class="fee-list">
        <div class="fee-card card" v-for="fee in fees" :key="fee.id">
          <div class="fee-header">
            <div class="fee-month">{{ fee.month }}</div>
            <span class="tag" :class="fee.status === 1 ? 'tag-success' : 'tag-warning'">
              {{ fee.status === 1 ? '已缴' : '未缴' }}
            </span>
          </div>
          
          <div class="fee-body">
            <div class="fee-amount">¥{{ fee.amount }}</div>
            <div class="fee-date" v-if="fee.pay_time">
              缴费时间: {{ formatDate(fee.pay_time) }}
            </div>
          </div>
          
          <div class="fee-footer" v-if="fee.status === 0">
            <button 
              class="btn btn-primary btn-sm" 
              @click="handlePay(fee)"
              :disabled="paying"
            >
              {{ paying ? '支付中...' : '立即缴费' }}
            </button>
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
      
      <div class="empty-state" v-if="fees.length === 0">
        <div class="empty-state-icon">💰</div>
        <p>暂无物业费账单</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getPropertyFeeList, payPropertyFee } from '@/api/service'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const fees = ref([])
const paying = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const fetchFees = async () => {
  try {
    const res = await getPropertyFeeList({
        page: currentPage.value,
        size: pageSize.value
    })
    fees.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取物业费列表失败:', error)
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchFees()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchFees()
}

const handlePay = async (fee) => {
  if (userStore.userInfo.balance < fee.amount) {
    ElMessage.warning('余额不足，请先充值')
    return
  }
  
  try {
    await ElMessageBox.confirm(`确认缴纳${fee.month}的物业费 ¥${fee.amount}？`, '缴费确认', {
      confirmButtonText: '立即缴费',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    paying.value = true
    await payPropertyFee({
      type: 2,  // 2表示物业费
      related_id: fee.id
    })
    ElMessage.success('缴费成功！')
    await fetchFees()
    await userStore.fetchUserInfo()  // 刷新余额
  } catch (error) {
    if (error !== 'cancel') {
        ElMessage.error('缴费失败：' + (error.response?.data?.msg || error.message))
    }
  } finally {
    paying.value = false
  }
}

onMounted(() => {
  fetchFees()
})
</script>

<style scoped>
.property-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.user-balance {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-lg);
}

.balance-amount {
  font-size: var(--font-size-2xl);
  font-weight: 700;
  color: var(--primary-color);
}

.fee-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.fee-card {
  padding: var(--spacing-lg);
}

.fee-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.fee-month {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--text-primary);
}

.fee-body {
  margin-bottom: var(--spacing-md);
}

.fee-amount {
  font-size: 32px;
  font-weight: 700;
  color: var(--danger-color);
  margin-bottom: var(--spacing-sm);
}

.fee-date {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.fee-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}
</style>
