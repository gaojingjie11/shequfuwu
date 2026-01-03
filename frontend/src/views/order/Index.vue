<template>
  <div class="order-page">
    <Navbar />
    
    <div class="container">
      <h1 class="page-title">我的订单</h1>
      
      <!-- Filter Tabs -->
      <div class="tabs-container">
          <el-tabs v-model="currentTab" @tab-click="handleTabClick">
            <el-tab-pane label="全部" name="all"></el-tab-pane>
            <el-tab-pane label="待支付" name="0"></el-tab-pane>
            <el-tab-pane label="待发货" name="1"></el-tab-pane>
            <el-tab-pane label="待收货" name="2"></el-tab-pane>
            <el-tab-pane label="已完成" name="3"></el-tab-pane>
          </el-tabs>
      </div>

      <div class="order-list" v-loading="loading">
        <template v-if="orders.length > 0">
           <el-card shadow="hover" class="order-item" v-for="order in orders" :key="order.id" :body-style="{ padding: '15px' }">
              <div slot="header" class="order-header">
                  <div class="header-left">
                      <span class="order-no">订单号: {{ order.order_no }}</span>
                      <span class="order-time">{{ formatDate(order.created_at) }}</span>
                  </div>
                  <el-tag :type="getStatusType(order.status)" size="small">
                      {{ getStatusText(order.status) }}
                  </el-tag>
              </div>

               <div class="order-products" @click="$router.push(`/order/${order.id}`)">
                  <div class="product-item" v-for="item in order.items" :key="item.id">
                      <el-image 
                        style="width: 60px; height: 60px; border-radius: 4px;"
                        :src="item.product?.image_url || 'https://via.placeholder.com/60'" 
                        fit="cover" />
                      <div class="product-info">
                          <div class="product-name">{{ item.product?.name || '未知商品' }}</div>
                          <div class="product-meta">
                              <span class="price">¥{{ item.price }}</span>
                              <span class="qty">x{{ item.quantity }}</span>
                          </div>
                      </div>
                  </div>
              </div>

              <div class="order-footer">
                  <div class="order-total">
                      共 {{ order.items.length }} 件商品  合计: <span class="total-price">¥{{ order.total_amount.toFixed(2) }}</span>
                  </div>
                  <div class="order-actions">
                       <el-button 
                          size="small" 
                          type="danger" 
                          v-if="order.status === 0"
                          @click.stop="payOrder(order.id)"
                          :loading="loadingMap[order.id]"
                        >
                          立即支付
                        </el-button>
                        <el-button 
                          size="small" 
                          v-if="order.status === 0"
                          @click.stop="cancelOrder(order.id)"
                          :loading="loadingMap[order.id]"
                        >
                          取消订单
                        </el-button>
                        <el-button 
                          size="small" 
                          type="success" 
                          v-if="order.status === 2"
                          @click.stop="confirmReceipt(order.id)"
                          :loading="loadingMap[order.id]"
                        >
                          确认收货
                        </el-button>
                  </div>
              </div>
           </el-card>

           <div class="pagination-container">
            <el-pagination
                v-model:current-page="page"
                v-model:page-size="size"
                :page-sizes="[10, 20, 50]"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="fetchOrders"
                @current-change="fetchOrders"
            />
            </div>
        </template>
        
        <el-empty v-else description="暂无相关订单" :image-size="200">
             <el-button type="primary" @click="$router.push('/mall')">去逛逛</el-button>
        </el-empty>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getOrderList, payOrder as payOrderApi, cancelOrder as cancelOrderApi, receiveOrder } from '@/api/order'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const orders = ref([])
const currentTab = ref('all')
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loadingMap = reactive({})

const handleTabClick = (pane) => {
    // Reset page to 1 when changing filters
    page.value = 1
    fetchOrders(pane.paneName)
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (status) => {
  const statusMap = {
    0: '待支付',
    1: '已支付',
    2: '已发货',
    3: '已完成',
    40: '已取消'
  }
  return statusMap[status] || '未知'
}

const getStatusType = (status) => {
  const map = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary',
    40: 'danger'
  }
  return map[status] || ''
}

const withLoading = async (id, fn) => {
  if (loadingMap[id]) return
  loadingMap[id] = true
  try {
    await fn()
  } finally {
    loadingMap[id] = false
  }
}

const payOrder = (orderId) => withLoading(orderId, async () => {
  try {
    await payOrderApi({ order_id: orderId })
    ElMessage.success('支付成功！')
    fetchOrders()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '支付失败')
  }
})

const confirmReceipt = (orderId) => withLoading(orderId, async () => {
  try {
    await receiveOrder(orderId)
    ElMessage.success('已确认收货')
    fetchOrders()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
})

const cancelOrder = (orderId) => withLoading(orderId, async () => {
  try {
    await ElMessageBox.confirm('确定要取消订单吗？', '提示', { type: 'warning' })
    await cancelOrderApi(orderId)
    ElMessage.success('已取消订单')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
        ElMessage.error('取消失败')
    }
  }
})

const fetchOrders = async (statusOverride) => {
  loading.value = true
  try {
    let statusParams = undefined
    
    // Check if statusOverride is valid (it might be an event object if called directly, so check type or value)
    // Actually, handleTabClick passes string/number.
    // If called from onMounted, it is undefined.
    // Use currentTab.value as fallback logic, but prefer override if it's a string/number
    
    let activeStatus = currentTab.value
    if (typeof statusOverride === 'string' || typeof statusOverride === 'number') {
        activeStatus = statusOverride
    }
    
    if (activeStatus !== 'all') {
        statusParams = parseInt(activeStatus)
    }

    const res = await getOrderList({ status: statusParams, page: page.value, size: size.value })
    if(res.list) {
        orders.value = res.list
        total.value = res.total
    } else {
        // Fallback for non-paginated API if used elsewhere
        orders.value = res
        total.value = res.length
    }
    
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
      loading.value = false
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.order-page {
  min-height: 100vh;
  padding-bottom: var(--spacing-xl);
}

.tabs-container {
    margin-bottom: 20px;
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.order-item {
    cursor: pointer;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 10px;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 15px;
}

.order-no {
  font-weight: 500;
  color: #333;
}

.order-time {
  color: #909399;
  font-size: 13px;
}

.order-products {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 10px 0;
}

.product-item {
  display: flex;
  gap: 15px;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #909399;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
  margin-top: 5px;
}

.order-total {
    font-size: 14px;
    color: #606266;
}

.total-price {
    color: #F56C6C;
    font-size: 16px;
    font-weight: 600;
    margin-left: 5px;
}

.order-actions {
  display: flex;
  gap: 10px;
}

.pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}
</style>
