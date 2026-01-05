<template>
  <div class="admin-child-page">
    <Navbar />
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">订单管理</h1>
        <div class="actions" style="display:flex; gap:10px;">
             <el-input 
                v-model="searchUserId" 
                placeholder="按用户ID搜索" 
                style="width: 200px" 
                clearable 
                @clear="handleSearch"
                @keyup.enter="handleSearch"
             >
                <template #append>
                    <el-button @click="handleSearch">搜索</el-button>
                </template>
             </el-input>
        </div>
      </div>

      <div class="table-container card">
        <el-table :data="orders" stripe border style="width: 100%" v-loading="loadingData">
          <el-table-column prop="order_no" label="订单号" width="180" />
          <el-table-column prop="user_id" label="用户ID" width="100" />
          <el-table-column prop="total_amount" label="总金额" width="120">
             <template #default="scope">¥{{ scope.row.total_amount }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
             <template #default="scope">
               <el-tag :type="getStatusType(scope.row.status)">
                 {{ getStatusText(scope.row.status) }}
               </el-tag>
             </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间">
             <template #default="scope">{{ formatDate(scope.row.created_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
             <template #default="scope">
               <el-button 
                 v-if="scope.row.status === 1" 
                 size="small" 
                 type="primary" 
                 @click="handleShip(scope.row)"
                 :loading="loadingMap[scope.row.id]"
               >
                 发货
               </el-button>
               <span v-else>-</span>
             </template>
          </el-table-column>
          <!-- Expandable row for order items -->
          <el-table-column type="expand">
              <template #default="scope">
                  <div style="padding: 10px 20px">
                      <p><strong>商品明细:</strong></p>
                      <el-table :data="scope.row.items" style="width: 100%" size="small">
                          <el-table-column label="商品图" width="80">
                              <template #default="itemScope">
                                  <img :src="itemScope.row.product?.image_url" style="width: 30px; height: 30px" />
                              </template>
                          </el-table-column>
                          <el-table-column prop="product.name" label="商品名称" />
                          <el-table-column prop="price" label="单价" />
                          <el-table-column prop="quantity" label="数量" />
                      </el-table>
                  </div>
              </template>
          </el-table-column>
        </el-table>

         <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="size"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="fetchOrders"
            @current-change="fetchOrders"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { getAdminOrderList, shipOrder } from '@/api/admin'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const orders = ref([])
const loadingMap = reactive({})
const loadingData = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)
const searchUserId = ref('')

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusText = (status) => {
  const map = { 0: '待支付', 1: '待发货', 2: '已发货', 3: '已完成', 40: '已取消' }
  return map[status] || status
}

const getStatusType = (status) => {
  if (status === 1) return 'warning' // Paid (Pending Ship) -> Yellow
  if (status === 2) return 'primary' // Shipped -> Blue
  if (status === 3) return 'success' // Completed -> Green
  if (status === 0) return 'danger'  // Unpaid -> Red
  if (status === 40) return 'info'   // Cancelled -> Grey
  return ''
}

const fetchOrders = async () => {
  loadingData.value = true
  try {
    const params = { page: page.value, size: size.value }
    if (searchUserId.value) {
        params.user_id = searchUserId.value
    }
    const res = await getAdminOrderList(params)
    // With simplified handler returning list, pagination might be missing if I didn't update API file
    // Handlers returns { list: [], total: N } now
    if (res.list) {
        orders.value = res.list
        total.value = res.total
    } else {
        // Fallback for old API style if something mismatched
        orders.value = res
        total.value = res.length
    }
  } catch (error) {
    console.error(error)
  } finally {
      loadingData.value = false
  }
}

const handleShip = async (order) => {
  if (loadingMap[order.id]) return
  loadingMap[order.id] = true
  
  try {
    await shipOrder({ id: order.id })
    ElMessage.success('发货成功')
    fetchOrders()
  } catch (error) {
    ElMessage.error('发货失败')
  } finally {
    loadingMap[order.id] = false
  }
}

const handleSearch = () => {
    page.value = 1
    fetchOrders()
}
 
onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.admin-child-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.table-container {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
