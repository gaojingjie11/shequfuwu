<template>
  <div class="transaction-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">我的账单</h1>
        <div class="balance-info">
            <el-statistic :value="userInfo.balance" :precision="2" prefix="当前余额: ¥" title="" value-style="color: var(--primary-color); font-weight: bold" />
        </div>
      </div>

      <el-card class="table-container">
        <el-table :data="transactions" stripe style="width: 100%" v-loading="loading">
            <el-table-column prop="created_at" label="时间" width="180">
                <template #default="scope">
                    {{ formatDate(scope.row.created_at) }}
                </template>
            </el-table-column>
            <el-table-column prop="type" label="类型" width="120">
                <template #default="scope">
                    <el-tag :type="getTypeTagType(scope.row.type)">{{ getTypeLabel(scope.row.type) }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="150">
                <template #default="scope">
                    <span :class="scope.row.amount > 0 ? 'text-success' : 'text-danger'">
                        {{ scope.row.amount > 0 ? '+' : '' }}{{ scope.row.amount.toFixed(2) }}
                    </span>
                </template>
            </el-table-column>
            <el-table-column prop="related_id" label="备注/单号" show-overflow-tooltip>
                <template #default="scope">
                     ID: {{ scope.row.related_id }} 
                     <span v-if="scope.row.id" class="text-gray"> (流水号: {{ scope.row.id }})</span>
                </template>
            </el-table-column>
        </el-table>
        
         <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchTransactions"
            @current-change="fetchTransactions"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getTransactionList } from '@/api/finance'
import { getUserInfo } from '@/api/user'
import dayjs from 'dayjs'

const transactions = ref([])
const userInfo = ref({ balance: 0 })
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)

const fetchTransactions = async () => {
    loading.value = true
    try {
        const res = await getTransactionList({ page: page.value, size: size.value })
        // Adapt if backend returns different structure
        if (res.list) {
            transactions.value = res.list
            total.value = res.total
        } else if (Array.isArray(res)) {
            // Fallback for older API structure
            transactions.value = res
            total.value = res.length
        }
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

const fetchUser = async () => {
    try {
        const res = await getUserInfo()
        userInfo.value = res
    } catch (e) {}
}

const formatDate = (dateStr) => {
    if(!dateStr) return ''
    return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const getTypeLabel = (type) => {
    const map = {
        1: '商城订单',
        2: '物业费',
        3: '充值',
        4: '转账'
    }
    return map[type] || '未知'
}

const getTypeTagType = (type) => {
    if (type === 3 || type === 4) return 'success'
    if (type === 1) return 'warning'
    if (type === 2) return 'info'
    return ''
}

onMounted(() => {
    fetchTransactions()
    fetchUser()
})
</script>

<style scoped>
.transaction-page {
    min-height: 100vh;
    padding-bottom: var(--spacing-xl);
}
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-lg);
}
.text-success { color: #67C23A; font-weight: bold; }
.text-danger { color: #F56C6C; font-weight: bold; }
.text-gray { color: #909399; font-size: 12px; }

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
