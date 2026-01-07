<template>
  <div class="order-detail-page">
    <Navbar />
    
    <div class="container">
      <div class="page-header">
        <el-page-header @back="$router.back()" content="订单详情" title="返回" />
      </div>

      <div class="content" v-loading="loading">
        <template v-if="order">
          <!-- Status Card -->
          <el-card class="status-card">
             <div class="status-header">
               <div class="status-left">
                 <span class="status-text">{{ getStatusText(order.status) }}</span>
                 <span class="store-name" v-if="order.store"> 
                    <el-tag size="small" type="info">服务门店</el-tag> {{ order.store.name }}
                 </span>
               </div>
               <span class="order-total">总金额: ¥{{ order.total_amount }}</span>
             </div>
             <div class="status-actions">
                <el-button v-if="order.status === 0" type="danger" @click="payOrder">立即支付</el-button>
                <el-button v-if="order.status === 0" @click="cancelOrder">取消订单</el-button>
                <el-button v-if="order.status === 2" type="primary" @click="confirmReceipt">确认收货</el-button>
             </div>
          </el-card>

          <!-- Address Info -->
          <el-card class="section-card">
            <template #header>
              <div class="card-header">
                <span>收货/用户信息</span>
              </div>
            </template>
            <!-- Data from SysUser association -->
            <div class="info-row" v-if="order.sys_user">
              <span class="label">用户名：</span>
              <span class="value">{{ order.sys_user.username }}</span>
            </div>
            <div class="info-row" v-if="order.sys_user">
               <span class="label">真实姓名：</span>
               <span class="value">{{ order.sys_user.real_name || '未设置' }}</span>
            </div>
            <div class="info-row" v-if="order.sys_user">
              <span class="label">手机号：</span>
              <span class="value">{{ order.sys_user.mobile }}</span>
            </div>
            <div class="info-row" v-else>
               <span class="value text-gray">用户信息获取失败</span>
            </div>
          </el-card>

          <!-- Order Info -->
          <el-card class="section-card">
             <template #header>
               <div class="card-header">
                 <span>订单信息</span>
               </div>
             </template>
             <div class="info-row">
               <span class="label">订单编号：</span>
               <span class="value">{{ order.order_no }}</span>
             </div>
             <div class="info-row">
               <span class="label">下单时间：</span>
               <span class="value">{{ formatDate(order.created_at) }}</span>
             </div>
          </el-card>

          <!-- Products -->
          <el-card class="section-card">
             <template #header>
               <div class="card-header">
                 <span>商品信息</span>
               </div>
             </template>
             <div class="product-list">
               <div v-for="item in order.items" :key="item.id" class="product-item" @click="$router.push(`/product/${item.product_id}`)">
                  <el-image 
                    :src="item.product_image || item.product?.image_url" 
                    class="product-img" 
                    fit="cover"
                  />
                  <div class="product-info">
                    <div class="product-name">{{ item.product_name || item.product?.name }}</div>
                    <div class="product-meta">
                      <span class="price">¥{{ item.price }}</span>
                      <span class="qty">x{{ item.quantity }}</span>
                    </div>
                  </div>
               </div>
             </div>
          </el-card>

        </template>
        <el-empty v-else description="未找到订单信息" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { getOrderDetail, payOrder as payOrderApi, cancelOrder as cancelOrderApi, receiveOrder } from '@/api/order'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const order = ref(null)
const loading = ref(false)

const getStatusText = (status) => {
  const map = {
    0: '待支付',
    1: '已支付',
    2: '已发货',
    3: '已完成',
    40: '已取消'
  }
  return map[status] || '未知状态'
}

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchDetail = async () => {
  loading.value = true
  try {
    const res = await getOrderDetail(route.params.id)
    order.value = res
  } catch (error) {
    ElMessage.error('获取订单详情失败')
  } finally {
    loading.value = false
  }
}

const payOrder = async () => {
  try {
    await payOrderApi({ order_id: order.value.id })
    ElMessage.success('支付成功')
    fetchDetail()
  } catch (e) {
    ElMessage.error('支付失败')
  }
}

const cancelOrder = async () => {
  try {
    await ElMessageBox.confirm('确定取消该订单吗?', '提示', { type: 'warning' })
    await cancelOrderApi(order.value.id)
    ElMessage.success('取消成功')
    fetchDetail()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败')
  }
}

const confirmReceipt = async () => {
   try {
    await receiveOrder(order.value.id)
    ElMessage.success('收货成功')
    fetchDetail()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  fetchDetail()
})
</script>

<style scoped>
.order-detail-page {
  min-height: 100vh;
  padding-bottom: 40px;
  background-color: #f5f7fa;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px 0;
}

.title {
  font-size: 18px;
  font-weight: 600;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.status-card {
  background: white;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.status-text {
  font-size: 24px;
  font-weight: bold;
  color: #E6A23C;
}

.order-total {
  font-size: 18px;
  color: #F56C6C;
  font-weight: bold;
}

.section-card {
  margin-bottom: 15px;
}

.info-row {
  display: flex;
  margin-bottom: 10px;
  font-size: 14px;
}

.label {
  color: #909399;
  width: 80px;
  flex-shrink: 0;
}

.value {
  color: #303133;
}

.product-item {
  display: flex;
  gap: 15px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}

.product-item:last-child {
  border-bottom: none;
}

.product-img {
  width: 80px;
  height: 80px;
  border-radius: 4px;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.product-name {
  font-size: 15px;
  margin-bottom: 8px;
  color: #333;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #666;
}
</style>
